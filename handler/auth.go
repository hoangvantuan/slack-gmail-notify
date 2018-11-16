package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/util"

	"github.com/mdshun/slack-gmail-notify/usecase"
	"golang.org/x/oauth2"

	"github.com/mdshun/slack-gmail-notify/infra"

	"github.com/labstack/echo"
)

const (
	slackAuthURL = "https://slack.com/oauth/authorize"
)

type authRequestParams struct {
	Code  string `query:"code" validate:"required"`
	State string `query:"state"`
}

type authHandler struct{}

// BindAuthHandler is handler for auth request
func BindAuthHandler(e *echo.Echo) {
	h := &authHandler{}

	e.GET("/v1/auth/slack", h.slackAuthURL)
	e.GET("/v1/auth/google", h.googleAuthURL)
	e.GET("/v1/auth/slack/redirected", h.slackAuth)
	e.GET("/v1/auth/google/redirected", h.googleAuth)
}

// redirect to auth slack page
func (a *authHandler) slackAuthURL(ctx echo.Context) error {
	url := fmt.Sprintf("%s?client_id=%s&scope=%s&redirect_uri=%s", slackAuthURL, infra.Env.SlackClientID, infra.Env.SlackScope, infra.Env.SlackRedirectedURL)

	infra.Sdebug("redirect to ", url)

	ctx.Redirect(http.StatusSeeOther, url)

	return nil
}

func (a *authHandler) googleAuthURL(ctx echo.Context) error {
	state := ctx.QueryParam("state")

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

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	infra.Sdebug("redirect to ", url)

	ctx.Redirect(http.StatusSeeOther, url)

	return nil
}

// get token for team
func (a *authHandler) slackAuth(ctx echo.Context) error {
	rp := &authRequestParams{}

	err := ctx.Bind(rp)
	if err == nil {
		err = ctx.Validate(rp)
	}

	if err != nil {
		infra.Swarn("error parameter is invalid", err)
		return ctx.String(http.StatusBadRequest, "error parameter is invalid")
	}

	uc := usecase.NewAuthUsecase()

	ri := &usecase.AuthRequestInput{
		Code:  rp.Code,
		State: rp.State,
	}

	err = uc.SlackAuth(ri)
	if err != nil {
		infra.Swarn("error while save database", err)
		return ctx.String(http.StatusInternalServerError, "error while save database")
	}

	return ctx.String(http.StatusOK, "thanks you for install app")
}

// TODO: need implements
func (a *authHandler) googleAuth(ctx echo.Context) error {
	rp := &authRequestParams{}

	state := ctx.QueryParam("state")
	decodedState, _ := util.Decrypt(state, infra.Env.EncryptKey)

	infra.Sdebug("authen google token for ", decodedState)

	secretInfo := &usecase.CommandRequestParams{}
	err := json.Unmarshal([]byte(decodedState), secretInfo)

	if err != nil {
		infra.Swarn("error parameter is invalid", err)
		return ctx.String(http.StatusBadRequest, "error parameter is invalid")
	}

	err = ctx.Bind(rp)
	if err == nil {
		err = ctx.Validate(rp)
		err = ctx.Validate(secretInfo)
	}

	if err != nil {
		infra.Swarn("error parameter is invalid", err)
		return ctx.String(http.StatusBadRequest, "error parameter is invalid")
	}

	uc := usecase.NewAuthUsecase()

	ri := &usecase.AuthRequestInput{
		Code:  rp.Code,
		State: rp.State,
	}

	err = uc.GoogleAuth(ri, secretInfo)

	if err != nil {
		infra.Swarn("error while save database", err)
		return ctx.String(http.StatusInternalServerError, "error while save database")
	}

	return ctx.String(http.StatusOK, "thank you, add gmail account successfull!")
}
