package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/nlopes/slack"
)

// SlackReqAuthMiddleware is verify slack request
func SlackReqAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// only auth for slack request
			if strings.Contains(c.Request().URL.Path, "/slack/") && !strings.Contains(c.Request().URL.Path, "/redirected") {
				var bodyBytes []byte
				if c.Request().Body != nil {
					bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
				}

				// Restore the io.ReadCloser to its original state
				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
				err := validateSlackReq(c.Request().Header, bodyBytes)
				if err != nil {
					return err
				}
			}

			return next(c)
		}
	}
}

func validateSlackReq(header http.Header, body []byte) error {
	sv, err := slack.NewSecretsVerifier(header, infra.Env.SlackSignSecret)
	if err != nil {
		return err
	}
	sv.Write(body)
	if err := sv.Ensure(); err != nil {
		return err
	}

	return nil
}
