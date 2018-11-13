package app

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	"github.com/mdshun/slack-gmail-notify/infra"

	"github.com/labstack/echo"
	"github.com/nlopes/slack"
)

var (
	// ErrCanNotVefiryRequest is can not vefiry error
	ErrCanNotVefiryRequest = errors.New("can not verify slack request")
)

// SlackReqAuthMiddleware is verify slack request
func SlackReqAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// only auth for slack request
			if strings.Contains(c.Request().URL.Path, "/slack") {
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
			} else {
				return next(c)
			}

			return next(c)
		}
	}
}

func validateSlackReq(header http.Header, body []byte) error {
	sv, err := slack.NewSecretsVerifier(header, infra.Env.SlackSignSecret)

	if err != nil {
		infra.Swarn(ErrCanNotVefiryRequest, err)
		return errors.Wrap(err, ErrCanNotVefiryRequest.Error())
	}

	sv.Write(body)

	if err := sv.Ensure(); err != nil {
		infra.Swarn(ErrCanNotVefiryRequest, err)
		return errors.Wrap(err, ErrCanNotVefiryRequest.Error())
	}

	return nil
}
