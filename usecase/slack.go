package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// slackAPI return slack client with given teamID
func slackAPI(teamID string) (*slack.Client, error) {
	teamRepo := rdb.NewTeamRepository(infra.RDB)

	team, err := teamRepo.FindByTeamID(teamID)
	if err != nil {
		return nil, errors.Wrap(err, "have error while find team")
	}

	token, err := util.Decrypt(team.BotAccessToken, infra.Env.EncryptKey)
	if err != nil {
		return nil, errors.Wrap(err, "error while decrypt token")
	}

	return slack.New(token), nil
}
