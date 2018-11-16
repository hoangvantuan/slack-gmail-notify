package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/usecase"

	"github.com/labstack/echo"
	"github.com/mdshun/slack-gmail-notify/infra"
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

	infra.Sdebug("payload ", payload)

	if err := json.Unmarshal([]byte(payload), rp); err != nil {
		infra.Swarn(errCanNotBindParam, err)
		return nil
	}

	infra.Sdebug("payload: ", rp)

	// Close button
	if rp.Actions[0].Name == "close" {
		return ctx.JSON(http.StatusOK, struct {
			ResponseType    string `json:"response_type"`
			Text            string `json:"text"`
			ReplaceOriginal bool   `json:"replace_original"`
			DeleteOriginal  bool   `json:"delete_original"`
		}{
			ResponseType:    "ephemeral",
			ReplaceOriginal: true,
			DeleteOriginal:  true,
		})
	}

	uc := usecase.NewIteractiveUsecase()

	// implements setting button
	if rp.Actions[0].Name == "setting" {
		uc.OpenSettingDialog(rp)
	}

	return ctx.NoContent(http.StatusOK)
}
