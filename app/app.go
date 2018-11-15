package app

import (
	"fmt"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/handler"
	"github.com/mediadotech/distribution-backend/cmd/public-api/validator"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mdshun/slack-gmail-notify/infra"
)

// Run is start app
func Run() {
	e := echo.New()

	// add validator
	e.Validator = validator.NewDefaultValidator()

	e.Debug = !infra.IsProduction()

	// Add verify slack request
	e.Use(SlackReqAuthMiddleware())

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"*"},
	}))

	handler.BindAuthHandler(e)
	handler.BindEventHandler(e)
	handler.BindCommandHandler(e)

	// Set API root.
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Gmail Slack Notify API server works!")
	})

	infra.Setup()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", infra.Env.Port)))
}
