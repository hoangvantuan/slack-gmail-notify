package util

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// SlackAPI return slack client with given teamID
func SlackAPI(teamID string) (*slack.Client, error) {
	teamRepo := rdb.NewTeamRepository(infra.RDB)

	team, err := teamRepo.FindByTeamID(teamID)
	if err != nil {
		return nil, errors.Wrap(err, "have error while find team")
	}

	token, err := Decrypt(team.BotAccessToken)
	if err != nil {
		return nil, errors.Wrap(err, "error while decrypt token")
	}

	return slack.New(token), nil
}
