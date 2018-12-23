package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
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
	SlackAuth(ri *AuthRequestInput) error
	GoogleAuth(ri *AuthRequestInput, rp *CommandRequestParams) error
}

// NewAuthUsecase will return auth usecase
func NewAuthUsecase() AuthUsecase {
	return &authUsecaseImpl{}
}

func (a *authUsecaseImpl) SlackAuth(ri *AuthRequestInput) error {
	or, err := slack.GetOAuthResponse(infra.Env.SlackClientID, infra.Env.SlackClientSecret, ri.Code, infra.Env.SlackRedirectedURL, infra.IsProduction())
	if err != nil {
		return errors.Wrap(err, "have error while get slack token")
	}

	tx := infra.RDB.Begin()

	team := &rdb.Team{}

	// encode token
	or.AccessToken, err = util.Encrypt(or.AccessToken, infra.Env.EncryptKey)
	or.Bot.BotAccessToken, err = util.Encrypt(or.Bot.BotAccessToken, infra.Env.EncryptKey)
	if err != nil {
		return errors.Wrap(err, "have error while encypt token")
	}

	// check team instated?
	team.AccessToken = or.AccessToken
	team.Scope = or.Scope
	team.TeamName = or.TeamName
	team.TeamID = or.TeamID
	team.UserID = or.UserID
	team.BotAccessToken = or.Bot.BotAccessToken
	team.BotUserID = or.Bot.BotUserID

	teamRepo := rdb.NewTeamRepository(tx)

	// check team was installed
	oldteam, err := teamRepo.FindByTeamID(team.TeamID)
	if err != nil {
		// is new team
		_, err = teamRepo.Add(team)
		if err != nil {
			return errors.Wrap(err, "have error while save team")
		}

		// save user
		user := &rdb.User{}
		user.UserID = team.UserID
		user.TeamID = team.TeamID

		userRepo := rdb.NewUserRepository(tx)

		_, err = userRepo.Add(user)
		if err != nil {
			tx.Rollback()
			return errors.Wrap(err, "have error while save user")
		}
	}

	// have old team
	if oldteam != nil {
		oldteam.AccessToken = or.AccessToken
		oldteam.Scope = or.Scope
		oldteam.TeamName = or.TeamName
		oldteam.TeamID = or.TeamID
		oldteam.UserID = or.UserID
		oldteam.BotAccessToken = or.Bot.BotAccessToken
		oldteam.BotUserID = or.Bot.BotUserID

		_, err = teamRepo.Update(oldteam)
		if err != nil {
			return errors.Wrap(err, "have error while save team")
		}
	}

	tx.Commit()

	return nil
}

func (a *authUsecaseImpl) GoogleAuth(ri *AuthRequestInput, rp *CommandRequestParams) error {
	ctx := context.Background()
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

	token, err := conf.Exchange(ctx, ri.Code)
	if err != nil {
		return errors.Wrap(err, "have error while get google token")
	}

	// get gmail
	client := conf.Client(ctx, token)

	srv, err := gmail.New(client)
	if err != nil {
		return errors.Wrap(err, "have error while creare gmail service")
	}

	gUserProfileCall := srv.Users.GetProfile("me")
	gUserProfileCall.Fields(googleapi.Field("emailAddress"))
	profile, err := gUserProfileCall.Do()
	if err != nil {
		return errors.Wrap(err, "have error while fetch list email")
	}

	// encode token
	token.AccessToken, err = util.Encrypt(token.AccessToken, infra.Env.EncryptKey)
	token.RefreshToken, err = util.Encrypt(token.RefreshToken, infra.Env.EncryptKey)
	if err != nil {
		return errors.Wrap(err, "have error while encypt token")
	}

	mygmail := &rdb.Gmail{
		UserID:       rp.UserID,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiryDate:   token.Expiry,
		TokenType:    token.TokenType,
		Email:        profile.EmailAddress,
	}

	gmailRepo := rdb.NewGmailRepository(infra.RDB)

	// check email was added
	oldgmail, err := gmailRepo.FindByEmail(profile.EmailAddress, rp.UserID)
	if err != nil {
		_, err = gmailRepo.Add(mygmail)
		if err != nil {
			return errors.Wrap(err, "have error while save email")
		}
	} else {
		oldgmail.AccessToken = token.AccessToken
		oldgmail.RefreshToken = token.RefreshToken
		oldgmail.ExpiryDate = token.Expiry
		oldgmail.TokenType = token.TokenType

		_, err = gmailRepo.Update(oldgmail)
		if err != nil {
			return errors.Wrap(err, "have error while update email")
		}
	}

	return nil
}
