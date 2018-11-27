package worker

import (
	"context"

	"github.com/carlescere/scheduler"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	gm "google.golang.org/api/gmail/v1"
)

const (
	labelUnread = "UNREAD"
)

// Job is define job with ID
type jobs map[int]*scheduler.Job

// Jobs is list job
var Jobs jobs

// Notifiers is interface defind method for job
type Notifiers interface {
	NotifyTeams() error
	NotifyTeam(team *rdb.Team) error
	NotifyUser(user *rdb.User, apiSlack *slack.Client) error
	NotifyGmail(gmail *rdb.Gmail, apiSlack *slack.Client) error
	StopNotifyTeam(team *rdb.Team) error
	StopNotifyUser(user *rdb.User, apiSlack *slack.Client) error
	StopNotifyGmail(gmail *rdb.Gmail, apiSlack *slack.Client) error
}

// SetupNotifiers is begin notify
func SetupNotifiers() error {
	Jobs = make(jobs)

	err := Jobs.NotifyTeams()
	if err != nil {
		panic(err)
	}

	return nil
}

func (j *jobs) NotifyTeams() error {
	teamRepo := rdb.NewTeamRepository(infra.RDB)

	teams, err := teamRepo.FindAllTeam()
	if err != nil {
		return errors.Wrap(err, "error while get all team")
	}

	for _, team := range teams {
		err = j.NotifyTeam(&team)
		if err != nil {
			return errors.Wrap(err, "error while notify all team")
		}
	}

	return nil
}

func (j *jobs) NotifyTeam(team *rdb.Team) error {
	slackAPI, err := util.SlackAPI(team.TeamID)
	if err != nil {
		return errors.Wrap(err, "error while init slack client")
	}

	userRepo := rdb.NewUserRepository(infra.RDB)
	users, err := userRepo.FindAllByTeamID(team.TeamID)
	if err != nil {
		return errors.Wrap(err, "error while get all users")
	}

	for _, user := range users {
		err := j.NotifyUser(&user, slackAPI)
		if err != nil {
			return errors.Wrap(err, "error while notify user")
		}
	}
	return nil
}

func (j *jobs) NotifyUser(user *rdb.User, apiSlack *slack.Client) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)

	gmails, err := gmailRepo.FindByUserID(user.UserID)
	if err != nil {
		return errors.Wrap(err, "error while get list gmails of user")
	}

	for _, gmail := range gmails {
		err = j.NotifyGmail(&gmail, apiSlack)
		if err != nil {
			return errors.Wrap(err, "error while notify gmail")
		}
	}
	return nil
}

func (j *jobs) NotifyGmail(gmail *rdb.Gmail, apiSlack *slack.Client) error {
	// check email is working in job
	j.StopNotifyGmail(gmail)

	job, err := scheduler.Every(5).Seconds().Run(func() {
		notify(gmail, apiSlack)
	})
	if err != nil {
		return errors.Wrap(err, "have error while notify email")
	}

	// add job
	(*j)[gmail.ID] = job

	return nil
}

func (j *jobs) StopNotifyTeam(team *rdb.Team) error {
	userRepo := rdb.NewUserRepository(infra.RDB)
	users, err := userRepo.FindAllByTeamID(team.TeamID)
	if err != nil {
		return errors.Wrap(err, "error while get all users")
	}

	for _, user := range users {
		err := j.StopNotifyUser(&user)
		if err != nil {
			return errors.Wrap(err, "error while notify user")
		}
	}
	return nil
}

func (j *jobs) StopNotifyUser(user *rdb.User) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)

	gmails, err := gmailRepo.FindByUserID(user.UserID)
	if err != nil {
		return errors.Wrap(err, "error while get list gmails of user")
	}

	for _, gmail := range gmails {
		err = j.StopNotifyGmail(&gmail)
		if err != nil {
			return errors.Wrap(err, "error while notify gmail")
		}
	}
	return nil
}

func (j *jobs) StopNotifyGmail(gmail *rdb.Gmail) error {
	if (*j)[gmail.ID] != nil {
		infra.Sdebug("stopping notify for email ", gmail.ID)
		curJob := (*j)[gmail.ID]
		curJob.Quit <- true
	}

	return nil
}

func notify(gmail *rdb.Gmail, apiSlack *slack.Client) {
	if gmail.NotifyChannelID == "" {
		return
	}

	accessToken, _ := util.Decrypt(gmail.AccessToken, infra.Env.EncryptKey)
	refreshToken, _ := util.Decrypt(gmail.RefreshToken, infra.Env.EncryptKey)

	conf := &oauth2.Config{
		ClientID:     infra.Env.GoogleClientID,
		ClientSecret: infra.Env.GoogleClientSecret,
		Scopes:       infra.Env.GoogleScopes,
		RedirectURL:  infra.Env.GoogleRedirectedURL,
		Endpoint: oauth2.Endpoint{
			AuthURL:  infra.Env.GoogleAuthURL,
			TokenURL: infra.Env.GoogleTokenURL,
		},
	}
	token := &oauth2.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    gmail.TokenType,
		Expiry:       gmail.ExpiryDate,
	}

	// get gmail
	client := conf.Client(context.Background(), token)

	srv, err := gm.New(client)
	if err != nil {
		infra.Sdebug(err, "have error while creare gmail service")
		return
	}

	msgRes, err := srv.Users.Messages.List("me").LabelIds(labelUnread).Do()
	if err != nil {
		infra.Sdebug("have error while get gmail", err)
		return
	}

	for _, msg := range msgRes.Messages {
		// Remove UNREAD label
		_, err := srv.Users.Messages.Modify("me", msg.Id, &gm.ModifyMessageRequest{
			RemoveLabelIds: []string{labelUnread},
		}).Do()
		if err != nil {
			infra.Sdebug("can not remove unread label ", msg.Id, err)
			return
		}

		_, _, err = apiSlack.PostMessage(gmail.NotifyChannelID, msg.Id, slack.PostMessageParameters{})
		if err != nil {
			infra.Sdebug("have error while post message", err)
			return
		}
	}
}
