package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// AuthRequestInput is auth request param
type AuthRequestInput struct {
	Code  string
	State string
}

type authUsecaseImpl struct{}

// AuthUsecase is auth interface
type AuthUsecase interface {
	SlackAuth(ri *AuthRequestInput) error
}

// NewAuthUsecase will return auth usecase
func NewAuthUsecase() AuthUsecase {
	return &authUsecaseImpl{}
}

func (a *authUsecaseImpl) SlackAuth(ri *AuthRequestInput) error {
	or, err := slack.GetOAuthResponse(infra.Env.SlackClientID, infra.Env.SlackClientSecret, ri.Code, infra.Env.SlackRedirectedURL, infra.IsProduction())
	if err != nil {
		infra.Swarn("has error while get token from code", ri, err)
		return errors.Wrap(err, "has error while get token from code")
	}

	tx := infra.RDB.Begin()

	team := &repository.Team{}

	infra.Sdebug("auth new team ", or)

	team.AccessToken = or.AccessToken
	team.Scope = or.Scope
	team.TeamName = or.TeamName
	team.TeamID = or.TeamID
	team.UserID = or.UserID
	team.BotAccessToken = or.Bot.BotAccessToken
	team.BotUserID = or.Bot.BotUserID

	infra.Sdebug("save team info ", team)

	teamRepo := repository.NewTeamRepository(tx)

	_, err = teamRepo.Add(team)
	if err != nil {
		tx.Rollback()
		infra.Swarn("have error while save team info", err)
		return errors.Wrap(err, "have error while save team info")
	}

	// save user
	user := &repository.User{}
	user.UserID = team.UserID
	user.TeamID = team.TeamID

	infra.Sdebug("save user info ", user)

	userRepo := repository.NewUserRepository(tx)

	_, err = userRepo.Add(user)
	if err != nil {
		tx.Rollback()
		infra.Swarn("have error while save user info ", err)
		return err
	}

	tx.Commit()

	return nil
}
