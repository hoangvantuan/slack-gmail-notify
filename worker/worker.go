package worker

import (
	"encoding/base64"
	"fmt"

	"github.com/carlescere/scheduler"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	labelUnread = "UNREAD"
	fetchTimes  = 10
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

	job, err := scheduler.Every(fetchTimes).Seconds().Run(func() {
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

	srv, err := util.GmailSrv(&oauth2.Token{
		AccessToken:  gmail.AccessToken,
		RefreshToken: gmail.RefreshToken,
		TokenType:    gmail.TokenType,
		Expiry:       gmail.ExpiryDate,
	})
	if err != nil {
		infra.Swarn(err, "have error while creare gmail service")
		return
	}

	msgRes, err := srv.Users.Messages.List("me").LabelIds(labelUnread).Do()
	if err != nil {
		infra.Swarn("have error while get gmail", err)
		return
	}

	var ids []string
	for _, msg := range msgRes.Messages {
		ids = append(ids, msg.Id)
		infra.Sdebug(msg.Id)

		msgDetails, err := srv.Users.Messages.Get("me", msg.Id).Do()
		if err != nil {
			infra.Swarn("can not get detail gmail", msg.Id)
			return
		}

		// get data
		var from, subject, cc, content string
		for _, header := range msgDetails.Payload.Headers {
			if header.Name == "From" {
				from = header.Value
			}

			if header.Name == "Subject" {
				subject = header.Value
			}

			if header.Name == "Cc" {
				cc = header.Value
			}
		}

		if msgDetails.Payload.MimeType == "text/plain" || msgDetails.Payload.MimeType == "multipart/alternative" {
			for _, part := range msgDetails.Payload.Parts {
				if part.MimeType == "text/plain" {
					contentByte, err := base64.URLEncoding.DecodeString(part.Body.Data)
					if err != nil {
						infra.Swarn("can not decode message", err)
						content = "can not parse message"
					} else {
						content = string(contentByte)
					}
				}
			}
		} else {
			content = msgDetails.Snippet
		}

		_, err = apiSlack.UploadFile(slack.FileUploadParameters{
			Filetype: "post",
			Channels: []string{gmail.NotifyChannelID},
			Content:  fmt.Sprintf("## FROM: %s\n## CC: %s\n\n%s", from, cc, content),
			Filename: subject,
		})
		if err != nil {
			infra.Swarn("have error while post message", err)
			return
		}
	}

	// Remove all UNREAD label
	// err = srv.Users.Messages.BatchModify("me", &gm.BatchModifyMessagesRequest{
	// 	Ids: ids,
	// }).Do()
	// if err != nil {
	// 	infra.Sdebug("can not remove unread label ", ids, err)
	// 	return
	// }
}
