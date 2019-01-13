package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/usecase"
	"github.com/mdshun/slack-gmail-notify/util"
	"golang.org/x/oauth2"
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
	//url := fmt.Sprintf("%s?client_id=%s&scope=%s&redirect_uri=%s", slackAuthURL, infra.Env.SlackClientID, infra.Env.SlackScope, infra.Env.SlackRedirectedURL)
	url := infra.SlackOauth2Config().AuthCodeURL("")
	ctx.Redirect(http.StatusSeeOther, url)
	return nil
}

func (a *authHandler) googleAuthURL(ctx echo.Context) error {
	state := ctx.QueryParam("state")
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := infra.GoogleOauth2Config().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	ctx.Redirect(http.StatusSeeOther, url)
	return nil
}

// get token for team
func (a *authHandler) slackAuth(ctx echo.Context) error {
	rp := &authRequestParams{}

	err := ctx.Bind(rp)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusBadRequest, "Parameter is invalid")
	}

	err = ctx.Validate(rp)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusBadRequest, "Parameter is invalid")
	}

	uc := usecase.NewAuthUsecase()
	err = uc.AuthSlack(&usecase.AuthRequestInput{
		Code:  rp.Code,
		State: rp.State,
	})
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusInternalServerError, "Slack authentication is fail")
	}

	return ctx.String(http.StatusOK, "thanks you for install app")
}

func (a *authHandler) googleAuth(ctx echo.Context) error {
	rp := &authRequestParams{}

	state := ctx.QueryParam("state")
	decodedState, _ := util.Decrypt(state)
	secretInfo := &usecase.UserIdentity{}
	err := json.Unmarshal([]byte(decodedState), secretInfo)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusBadRequest, "Parameter is invalid")
	}
	err = ctx.Bind(rp)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusBadRequest, "Parameter is invalid")
	}

	err = ctx.Validate(rp)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusBadRequest, "Parameter is invalid")
	}

	err = ctx.Validate(secretInfo)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusBadRequest, "Parameter is invalid")
	}

	uc := usecase.NewAuthUsecase()
	err = uc.AuthGoogle(&usecase.AuthRequestInput{
		Code:  rp.Code,
		State: rp.State,
	}, secretInfo)
	if err != nil {
		infra.Warn(err)
		return ctx.String(http.StatusInternalServerError, "Google account authentication is fail")
	}

	return ctx.String(http.StatusOK, "thanks you, register gmail account successfull!")
}
