package handler

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/usecase"
	"github.com/nlopes/slack"
)

type iteractiveHandler struct{}

// BindIteractiveHandler is handler for event
func BindIteractiveHandler(e *echo.Echo) {
	h := &iteractiveHandler{}

	e.POST("/v1/slack/interactive", h.handler)
}

func (e *iteractiveHandler) handler(ctx echo.Context) error {
	rp := &usecase.IteractiveRequestParams{}

	payload := ctx.FormValue("payload")

	if err := json.Unmarshal([]byte(payload), rp); err != nil {
		return ctx.NoContent(http.StatusOK)
	}

	// Close button
	if rp.Actions[0].Name == "close" {
		return ctx.JSON(http.StatusOK, slack.Msg{
			ResponseType:    "ephemeral",
			ReplaceOriginal: true,
			DeleteOriginal:  true,
		})
	}

	uc := usecase.NewIteractiveUsecase()

	if rp.Actions[0].Name == "list-gmail" {
		err := uc.ListAccount(rp)
		if err != nil {
			return ctx.NoContent(http.StatusOK)
		}
	}

	if rp.Actions[0].Name == "notify-channel" {
		err := uc.NotifyChannel(rp)
		if err != nil {
			return ctx.NoContent(http.StatusOK)
		}
	}

	if rp.Actions[0].Name == "remove-gmail" {
		err := uc.RemoveAccount(rp)
		if err != nil {
			return ctx.NoContent(http.StatusOK)
		}
	}

	return ctx.NoContent(http.StatusOK)
}
