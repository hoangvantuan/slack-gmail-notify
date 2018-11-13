package handler

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/pkg/errors"

	"github.com/labstack/echo"
	"github.com/nlopes/slack/slackevents"
)

type verificationRequestParam struct {
	Token     string `json:"token" form:"token"`
	Challenge string `json:"challenge" form:"challenge"`
	Type      string `json:"type" form:"type"`
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
		infra.Swarn("has error while parse event request ", err)
		return errors.Wrap(err, "has error while parse event requests")
	}

	if eventAPI.Type == slackevents.URLVerification {
		return verificationEventHandler(ctx)
	}

	return nil
}

func verificationEventHandler(ctx echo.Context) error {
	vr := &verificationRequestParam{}

	if err := ctx.Bind(vr); err != nil {
		infra.Swarn("can not binding challenge request ", err)
		return errors.Wrap(err, "can not binding challenge request")
	}

	ctx.String(http.StatusOK, vr.Challenge)

	return nil
}
