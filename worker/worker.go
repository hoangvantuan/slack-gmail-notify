package worker

import (
	"github.com/carlescere/scheduler"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
)

const (
	fetchTimes = 10
)

// Job is define job with ID
var jobs map[int]*scheduler.Job

type messages struct {
	m   []*message
	ids []string
}

type message struct {
	ID      string
	From    string
	CC      string
	Subject string
	Body    string
}

// Setup -
func Setup() {
	jobs = make(map[int]*scheduler.Job)

	err := NotifyTeams()
	if err != nil {
		panic(err)
	}
}

// NotifyTeams -
func NotifyTeams() error {
	teamRepo := rdb.NewTeamRepository(infra.RDB)

	teams, err := teamRepo.FindAllTeam()
	if err != nil {
		return errors.Wrap(err, "error while get all team")
	}

	for _, team := range teams {
		err = NotifyTeam(&team)
		if err != nil {
			return errors.Wrap(err, "error while notify all team")
		}
	}

	return nil
}

// NotifyTeam -
func NotifyTeam(team *rdb.Team) error {
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
		err := NotifyUser(&user, slackAPI)
		if err != nil {
			return errors.Wrap(err, "error while notify user")
		}
	}
	return nil
}

// NotifyUser -
func NotifyUser(user *rdb.User, apiSlack *slack.Client) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)

	gmails, err := gmailRepo.FindByUserID(user.UserID)
	if err != nil {
		return errors.Wrap(err, "error while get list gmails of user")
	}

	for _, gmail := range gmails {
		err = NotifyGmail(&gmail, apiSlack)
		if err != nil {
			return errors.Wrap(err, "error while notify gmail")
		}
	}
	return nil
}

// NotifyGmail -
func NotifyGmail(gmail *rdb.Gmail, apiSlack *slack.Client) error {
	// check email is working in job
	StopNotifyGmail(gmail)

	job, err := scheduler.Every(fetchTimes).Seconds().Run(func() {
		notify(gmail, apiSlack)
	})
	if err != nil {
		return errors.Wrap(err, "have error while notify email")
	}

	// add job
	jobs[gmail.ID] = job

	return nil
}

// StopNotifyTeam -
func StopNotifyTeam(team *rdb.Team) error {
	userRepo := rdb.NewUserRepository(infra.RDB)
	users, err := userRepo.FindAllByTeamID(team.TeamID)
	if err != nil {
		return errors.Wrap(err, "error while get all users")
	}

	for _, user := range users {
		err := StopNotifyUser(&user)
		if err != nil {
			return errors.Wrap(err, "error while notify user")
		}
	}
	return nil
}

// StopNotifyUser -
func StopNotifyUser(user *rdb.User) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)

	gmails, err := gmailRepo.FindByUserID(user.UserID)
	if err != nil {
		return errors.Wrap(err, "error while get list gmails of user")
	}

	for _, gmail := range gmails {
		err = StopNotifyGmail(&gmail)
		if err != nil {
			return errors.Wrap(err, "error while notify gmail")
		}
	}
	return nil
}

// StopNotifyGmail -
func StopNotifyGmail(gmail *rdb.Gmail) error {
	if jobs[gmail.ID] != nil {
		curJob := jobs[gmail.ID]
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
		return
	}

	gw := newGGWorker(srv)
	ms, err := gw.fetchUnread()
	if err != nil {
		return
	}

	sw := newSlWorker(apiSlack)
	err = sw.posts(ms.m, gmail.NotifyChannelID)
	if err != nil {
		return
	}

	err = gw.read(ms)
	if err != nil {
		return
	}
}
