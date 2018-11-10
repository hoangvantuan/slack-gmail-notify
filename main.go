package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/mdshun/slack-gmail-notify/infra"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{"*"},
	}))

	// Set API root.
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Gmail Slack Notify API server works!")
	})

	infra.Setup()
	e.Logger.Fatal(e.Start(":80"))
}
