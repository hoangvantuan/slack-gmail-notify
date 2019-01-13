package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/usecase"
	"github.com/nlopes/slack/slackevents"
)

type verificationRequestParam struct {
	Token     string `json:"token" form:"token"`
	Challenge string `json:"challenge" form:"challenge"`
	Type      string `json:"type" form:"type"`
}

type appUninstallRequestParam struct {
	TeamID string `json:"team_id" form:"team_id"`
}

type eventHandler struct{}

// BindEventHandler is handler for event
func BindEventHandler(e *echo.Echo) {
	h := &eventHandler{}

	e.POST("/v1/slack/event", h.handler)
}

// TODO: need verification request
func (e *eventHandler) handler(ctx echo.Context) error {
	var bodyBytes []byte
	if ctx.Request().Body != nil {
		bodyBytes, _ = ioutil.ReadAll(ctx.Request().Body)
		// Restore the io.ReadCloser to its original state
		ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// disable verifytoken
	// TODO: use function of library
	optionNoVerifyToken := func(cfg *slackevents.Config) {
		cfg.TokenVerified = true
	}

	eventAPI, err := slackevents.ParseEvent(json.RawMessage(bodyBytes), optionNoVerifyToken)
	if err != nil {
		infra.Warn(err)
		return ctx.NoContent(http.StatusOK)
	}

	// handler verify url event
	if eventAPI.Type == slackevents.URLVerification {
		return verificationEventHandler(ctx)
	}

	if eventAPI.Type == slackevents.CallbackEvent {
		innerEvent := eventAPI.InnerEvent
		switch innerEvent.Data.(type) {
		// handler app uninstall event
		case *slackevents.AppUninstalledEvent:
			return appUninstall(ctx, eventAPI.TeamID)
		}
	}

	return ctx.NoContent(http.StatusOK)
}

func verificationEventHandler(ctx echo.Context) error {
	vr := &verificationRequestParam{}

	if err := ctx.Bind(vr); err != nil {
		infra.Warn(err)
		return ctx.NoContent(http.StatusOK)
	}

	return ctx.String(http.StatusOK, vr.Challenge)
}

func appUninstall(ctx echo.Context, teamID string) error {
	ctx.String(http.StatusOK, "ok")

	uc := usecase.NewEventUsecase()

	err := uc.UninstallApp(teamID)
	if err != nil {
		infra.Warn(err)
		return ctx.NoContent(http.StatusOK)
	}

	return ctx.NoContent(http.StatusOK)
}
