package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/nlopes/slack"
	gmail "google.golang.org/api/gmail/v1"
	"google.golang.org/api/googleapi"
)

// AuthRequestInput is auth request param
type AuthRequestInput struct {
	Code  string
	State string
}

type authUsecaseImpl struct{}

// AuthUsecase is auth interface
type AuthUsecase interface {
	AuthSlack(ri *AuthRequestInput) error
	AuthGoogle(ri *AuthRequestInput, rp *slack.SlashCommand) error
}

// NewAuthUsecase will return auth usecase
func NewAuthUsecase() AuthUsecase {
	return &authUsecaseImpl{}
}

func (a *authUsecaseImpl) AuthSlack(ri *AuthRequestInput) error {
	// Get authentication response from slack
	or, err := slack.GetOAuthResponse(infra.Env.SlackClientID, infra.Env.SlackClientSecret, ri.Code, infra.Env.SlackRedirectedURL, infra.IsProduction())
	if err != nil {
		return err
	}

	// start transaction
	tx := infra.RDB.Begin()
	teamRepo := rdb.NewTeamRepository(tx)
	err := teamRepo.Save(&rdb.Team{
		AccessToken:    or.AccessToken,
		Scope:          or.Scope,
		TeamName:       or.TeamName,
		TeamID:         or.TeamID,
		UserID:         or.UserID,
		BotAccessToken: or.Bot.BotAccessToken,
		BotUserID:      or.Bot.BotUserID,
	})
	if err != inl {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (a *authUsecaseImpl) AuthGoogle(ri *AuthRequestInput, rp *slack.SlashCommand) error {
	token, err := infra.GoogleOauth2Config().Exchange(ri.Code)
	if err != nil {
		return err
	}

	// get gmail
	client := infra.GoogleOauth2Config().Client(token)
	srv, err := gmail.New(client)
	if err != nil {
		return err
	}

	gUserProfileCall := srv.Users.GetProfile("me")
	gUserProfileCall.Fields(googleapi.Field("emailAddress"))
	profile, err := gUserProfileCall.Do()
	if err != nil {
		return err
	}

	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	err = gmailRepo.Save(&rdb.Gmail{
		UserID:       rp.UserID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiryDate:   token.Expiry,
		TokenType:    token.TokenType,
		Email:        profile.EmailAddress,
	})

	return err
}
