package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/carlescere/scheduler"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/nlopes/slack"
	"golang.org/x/oauth2"
	gm "google.golang.org/api/gmail/v1"
)

const (
	fetchTimes = 5
	running    = true
	stop       = false
)

type safeJobStatus struct {
	jobStatus map[string]bool
	mux       sync.Mutex
}

type safeJobs struct {
	jobs map[string]*scheduler.Job
	mux  sync.Mutex
}

func (s *safeJobStatus) set(k string, v bool) {
	s.mux.Lock()
	s.jobStatus[k] = v
	s.mux.Unlock()
}

func (s *safeJobStatus) get(k string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.jobStatus[k]
}

func (s *safeJobs) set(k string, v *scheduler.Job) {
	s.mux.Lock()
	s.jobs[k] = v
	s.mux.Unlock()
}

func (s *safeJobs) get(k string) *scheduler.Job {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.jobs[k]
}

func (s *safeJobs) has(k string) bool {
	s.mux.Lock()
	defer s.mux.Unlock()
	_, found := s.jobs[k]
	return found
}

func (s *safeJobs) delete(k string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	delete(s.jobs, k)
}

// Job is define job with ID
var jobs *safeJobs
var jobStatus *safeJobStatus

type messages struct {
	m   []*message
	ids []string
}

type message struct {
	ID      string
	From    string
	To      string
	CC      string
	Subject string
	Body    string
}

// Setup -
func Setup() {
	infra.Info("Setup worker...")
	jobs = &safeJobs{
		jobs: make(map[string]*scheduler.Job),
	}
	jobStatus = &safeJobStatus{
		jobStatus: make(map[string]bool),
	}

	if err := NotifyForTeams(); err != nil {
		panic(err)
	}

	go func() {
		for {
			infra.Debug(fmt.Sprintf("Have %d jobs is running ", len(jobs.jobs)))
			time.Sleep(time.Second * 60)
		}
	}()
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
	if gmail.Status == rdb.Stop {
		infra.Debug(gmail.Email, " is stop")
		return nil
	}

	if gmail.Status == rdb.Pending {
		infra.Debug(gmail.Email, " is pending")
		return nil
	}
	// check email is working in job
	StopNotifyForGmail(gmail.Email)
	infra.Debug("Start notify for ", gmail.Email)

	job, err := scheduler.Every(fetchTimes).Seconds().Run(func() {
		// not run if job is running
		if jobStatus.get(gmail.Email) {
			return
		}

		// set to running
		jobStatus.set(gmail.Email, running)

		err := notify(gmail, apiSlack)
		if err != nil {
			jobStatus.set(gmail.Email, stop)
			infra.Warn(err)
		}
	})
	if err != nil {
		return err
	}

	// add job
	jobs.set(gmail.Email, job)

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
	if jobs.has(mail) {
		infra.Debug("Stop notify for ", mail)
		curJob := jobs.get(mail)
		curJob.Quit <- true
		jobs.delete(mail)
		jobStatus.set(mail, stop)
	}
}

func notify(gmail *rdb.Gmail, apiSlack *slack.Client) error {
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

	infra.Info(fmt.Sprintf("Fetching new message for email (%s)", gmail.Email))
	gw := newGGWorker(srv, gmail.Email)
	ms, err := gw.fetchUnread()
	if err != nil {
		return err
	}

	if len(ms.ids) > 0 {
		infra.Info("Email ", gmail.Email, " has (", len(ms.ids), ") new message")
	}

	sw := newSlWorker(apiSlack)
	isRead := true
	if gmail.MarkAs == "" || gmail.MarkAs == "unread" {
		isRead = false
	}

	err = sw.posts(gw, ms.m, gmail.NotifyChannelID, isRead, gmail.LabelID)
	if err != nil {
		return err
	}

	// set to stop
	jobStatus.set(gmail.Email, stop)

	return nil
}
