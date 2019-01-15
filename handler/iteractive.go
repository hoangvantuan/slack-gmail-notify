package handler

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/usecase"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
)

type iteractiveHandler struct{}

// BindIteractiveHandler is handler for event
func BindIteractiveHandler(e *echo.Echo) {
	h := &iteractiveHandler{}

	e.POST("/v1/slack/interactive", h.withNoContent)
}

func (e *iteractiveHandler) handler(ctx echo.Context) {
	rp := &usecase.IteractiveRequestParams{}

	payload := ctx.FormValue("payload")

	if err := json.Unmarshal([]byte(payload), rp); err != nil {
		infra.Warn(err)
		return
	}

	if rp.Actions[0].Name == util.CloseName {
		res := slack.Msg{
			ResponseType:    "ephemeral",
			ReplaceOriginal: true,
			DeleteOriginal:  true,
		}
		resJSON, _ := json.Marshal(&res)
		_, err := http.Post(rp.ResponseURL, "application/json", bytes.NewReader(resJSON))
		if err != nil {
			infra.Warn(err)
		}
	}

	uc := usecase.NewIteractiveUsecase()

	if rp.Actions[0].Name == util.ListGmailAccountName {
		err := uc.ListAllAccount(rp)
		if err != nil {
			infra.Warn(err)
			return
		}
	}

	if rp.Actions[0].Name == util.NotifyChannelName {
		err := uc.NotifyToChannel(rp)
		if err != nil {
			infra.Warn(err)
			return
		}
	}

	if rp.Actions[0].Name == util.MarkAsName {
		err := uc.MarkAs(rp)
		if err != nil {
			infra.Warn(err)
			return
		}
	}

	if rp.Actions[0].Name == util.StartEmailName {
		err := uc.Start(rp)
		if err != nil {
			infra.Warn(err)
			return
		}
	}

	if rp.Actions[0].Name == util.StopEmailName {
		err := uc.Stop(rp)
		if err != nil {
			infra.Warn(err)
			return
		}
	}

	if rp.Actions[0].Name == util.RemmoveGmailAccountName {
		err := uc.RemoveAccount(rp)
		if err != nil {
			infra.Warn(err)
			return
		}
	}
}

func (e *iteractiveHandler) withNoContent(ctx echo.Context) error {
	return withNoContent(ctx, e.handler)
}
