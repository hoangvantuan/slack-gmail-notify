package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

func withNoContent(ctx echo.Context, fn func(ctx echo.Context)) error {
	go fn(ctx)
	return ctx.NoContent(http.StatusOK)
}
