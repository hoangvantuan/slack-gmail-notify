package handler

import (
	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/usecase"
	"github.com/nlopes/slack"
)

type commandHandler struct{}

// BindCommandHandler is handler for event
func BindCommandHandler(e *echo.Echo) {
	h := &commandHandler{}

	e.POST("/v1/slack/command", h.withNoContent)
}

func (e *commandHandler) handler(ctx echo.Context) {
	rp, err := slack.SlashCommandParse(ctx.Request())
	if err != nil {
		infra.Warn(err)
		return
	}
	uc := usecase.NewCommandUsecase()
	err = uc.GetMainMenu(&rp)
	if err != nil {
		infra.Warn(err)
	}
}

func (e *commandHandler) withNoContent(ctx echo.Context) error {
	return withNoContent(ctx, e.handler)
}
