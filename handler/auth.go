package handler

import (
	"fmt"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/usecase"

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

// get token for team
func (a *authHandler) slackAuth(ctx echo.Context) error {
	rp := &authRequestParams{}

	err := ctx.Bind(rp)
	if err == nil {
		err = ctx.Validate(rp)
	}

	if err != nil {
		infra.Swarn("request parameter is not valid", err)
		return ctx.String(http.StatusBadRequest, "request parameter is not valid")
	}

	uc := usecase.NewAuthUsecase()

	ri := &usecase.AuthRequestInput{
		Code:  rp.Code,
		State: rp.State,
	}

	err = uc.SlackAuth(ri)
	if err != nil {
		infra.Swarn("has error while save to database", err)
		return ctx.String(http.StatusInternalServerError, "has error with database")
	}

	return ctx.String(http.StatusOK, "thanks you for install app")
}

// TODO: need implements
func (a *authHandler) googleAuth(ctx echo.Context) error {
	return nil
}
