package handler

import (
	"net/http"

	"github.com/nlopes/slack"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/usecase"
)

type commandHandler struct{}

// BindCommandHandler is handler for event
func BindCommandHandler(e *echo.Echo) {
	h := &commandHandler{}

	e.POST("/v1/slack/command", h.handler)
}

func (e *commandHandler) handler(ctx echo.Context) (err error) {
	// always return 200 status
	defer func() {
		err = ctx.NoContent(http.StatusOK)
	}()

	rp := &slack.SlashCommand{}

	if err = ctx.Bind(rp); err != nil {
		infra.Warn(err)
		return
	}

	uc := usecase.NewCommandUsecase()
	err = uc.GetMainMenu(rp)
	if err != nil {
		infra.Warn(err)
		return
	}

	return
}
