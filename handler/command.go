package handler

import (
	"net/http"

	"github.com/mdshun/slack-gmail-notify/usecase"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
)

type commandHandler struct{}

// BindCommandHandler is handler for event
func BindCommandHandler(e *echo.Echo) {
	h := &commandHandler{}

	e.POST("/v1/slack/command", h.handler)
}

func (e *commandHandler) handler(ctx echo.Context) error {
	rp := &usecase.CommandRequestParams{}

	if err := ctx.Bind(rp); err != nil {
		infra.Swarn("error can not bind parameter", err)
		return ctx.NoContent(http.StatusOK)
	}

	uc := usecase.NewCommandUsecase()

	err := uc.MainMenu(rp)

	if err != nil {
		infra.Swarn("error while handler command", err)
		return ctx.NoContent(http.StatusOK)
	}

	return ctx.NoContent(http.StatusOK)

}
