package worker

import (
	"context"

	"github.com/carlescere/scheduler"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	gm "google.golang.org/api/gmail/v1"
)

const (
	fetchTimes = 2
)

// Job is define job with ID
var jobs map[string]*scheduler.Job
var jobStatus map[string]bool

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
	jobs = make(map[string]*scheduler.Job)
	jobStatus = make(map[string]bool)

	if err := NotifyForTeams(); err != nil {
		panic(err)
	}
}

// NotifyForTeams to all team
func NotifyForTeams() error {
	teamRepo := rdb.NewTeamRepository(infra.RDB)

	teams, err := teamRepo.FindAllTeam()
	if err != nil {
		return err
	}

	for _, team := range teams {
		err = NotifyForTeam(team)
		if err != nil {
			return err
		}
	}

	return nil
}

// NotifyForTeam to team
func NotifyForTeam(team *rdb.Team) error {
	infra.Debug("Starting notify for team ", team.TeamName)

	apiSlack := slack.New(team.BotAccessToken)
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	gmails, err := gmailRepo.FindByTeamID(team.TeamID)
	if err != nil {
		return err
	}

	for _, gmail := range gmails {
		err = NotifyForGmail(gmail, apiSlack)
		if err != nil {
			return err
		}
	}

	return nil
}

// NotifyForUser -
func NotifyForUser(user *rdb.User) error {
	teamRepo := rdb.NewTeamRepository(infra.RDB)
	team, err := teamRepo.FindByTeamID(user.TeamID)
	if err != nil {
		return err
	}

	apiSlack := slack.New(team.BotAccessToken)

	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	gmails, err := gmailRepo.FindByUser(user)
	if err != nil {
		return err
	}

	for _, gmail := range gmails {
		err = NotifyForGmail(gmail, apiSlack)
		if err != nil {
			return err
		}
	}
	return nil
}

// NotifyForGmail -
func NotifyForGmail(gmail *rdb.Gmail, apiSlack *slack.Client) error {
	infra.Debug("Starting notify for ", gmail.Email)

	// check email is working in job
	StopNotifyForGmail(gmail.Email)

	job, err := scheduler.Every(fetchTimes).Seconds().Run(func() {
		// not run if job is running
		if jobStatus[gmail.Email] {
			return
		}

		// set to running
		jobStatus[gmail.Email] = true

		err := notify(gmail, apiSlack)
		if err != nil {
			jobStatus[gmail.Email] = false
			infra.Warn(err)
		}
	})
	if err != nil {
		return err
	}

	// add job
	jobs[gmail.Email] = job

	return nil
}

// StopNotifyForTeam -
func StopNotifyForTeam(team *rdb.Team) error {
	infra.Debug("Stop notify for ", team.TeamName)

	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	gmails, err := gmailRepo.FindByTeamID(team.TeamID)
	if err != nil {
		return err
	}

	for _, gmail := range gmails {
		StopNotifyForGmail(gmail.Email)
	}

	return nil
}

// StopNotifyForUser -
func StopNotifyForUser(user *rdb.User) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	gmails, err := gmailRepo.FindByUser(user)
	if err != nil {
		return err
	}

	for _, gmail := range gmails {
		StopNotifyForGmail(gmail.Email)
	}

	return nil
}

// StopNotifyForGmail -
func StopNotifyForGmail(mail string) {
	if _, found := jobs[mail]; found {
		infra.Debug("Stop notify for ", mail)
		curJob := jobs[mail]
		curJob.Quit <- true
	}
}

func notify(gmail *rdb.Gmail, apiSlack *slack.Client) error {
	infra.Debug("Starting notify for ", gmail.Email)
	if gmail.NotifyChannelID == "" {
		return nil
	}

	// get gmail
	client := infra.GoogleOauth2Config().Client(context.Background(), &oauth2.Token{
		AccessToken:  gmail.AccessToken,
		RefreshToken: gmail.RefreshToken,
		TokenType:    gmail.TokenType,
		Expiry:       gmail.ExpiryDate,
	})

	srv, err := gm.New(client)
	if err != nil {
		return err
	}

	gw := newGGWorker(srv)
	ms, err := gw.fetchUnread()
	if err != nil {
		return err
	}

	infra.Debug("Email ", gmail.Email, " has (", len(ms.ids), ") new message")

	sw := newSlWorker(apiSlack)
	err = sw.posts(gw, ms.m, gmail.NotifyChannelID)
	if err != nil {
		return err
	}

	// set to stop
	jobStatus[gmail.Email] = false
	infra.Debug("Stop notify for ", gmail.Email)

	return nil
}
