package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mdshun/slack-gmail-notify/handler"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/worker"
)

// Run is start app
func Run() {
	e := echo.New()

	// add validator
	e.Validator = NewDefaultValidator()

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
	handler.BindIteractiveHandler(e)

	// Set API root.
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Gmail Slack Notify API server works!")
	})

	infra.Setup()
	infra.Info("Starting app in ", time.Now().Format("2006/1/2 15:04:05 MST"))
	worker.Setup()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", infra.Env.Port)))
}
