package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/util"

	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// CommandRequestParams is request from command
type CommandRequestParams struct {
	Token       string `json:"token" form:"token"`
	TeamID      string `json:"team_id" form:"team_id"`
	TeamDomain  string `json:"team_domain" form:"team_domain"`
	ChannelID   string `json:"channel_id" form:"channel_id"`
	ChannelName string `json:"channel_name" form:"channel_name"`
	UserID      string `json:"user_id" form:"user_id"`
	UserName    string `json:"user_name" form:"user_name"`
	Command     string `json:"command" form:"command"`
	Text        string `json:"text" form:"text"`
	ResponseURL string `json:"response_url" form:"response_url"`
	TriggerID   string `json:"trigger_id" form:"trigger_id"`
}

type commandUsecaseImpl struct{}

// CommandUsecase is event interface
type CommandUsecase interface {
	MainMenu(rp *CommandRequestParams) error
}

// NewCommandUsecase will return event usecase
func NewCommandUsecase() CommandUsecase {
	return &commandUsecaseImpl{}
}

func (c *commandUsecaseImpl) MainMenu(rp *CommandRequestParams) error {
	msgAt := genInteractiveMenu(rp)

	msgatstr, err := json.Marshal(msgAt)
	if err != nil {
		return errors.Wrap(err, "have error while decode json")
	}

	infra.Sdebug(string(msgatstr))

	_, err = http.Post(rp.ResponseURL, "application/json", bytes.NewReader(msgatstr))
	if err != nil {
		return errors.Wrap(err, "have error while post message")
	}

	return nil
}

func genInteractiveMenu(rp *CommandRequestParams) slack.Msg {
	at := slack.Attachment{}

	pjson, _ := json.Marshal(rp)
	pjsoneconded, _ := util.Encrypt(string(pjson), infra.Env.EncryptKey)
	// generate action
	addBtn := slack.AttachmentAction{
		Name:  "add-gmail-account",
		Text:  "Add Account",
		Value: "AddGmailAccount",
		Style: "primary",
		Type:  "button",
		URL:   fmt.Sprintf("%s/v1/auth/google?state=%s", infra.Env.APIHost, pjsoneconded),
	}

	settingBtn := slack.AttachmentAction{
		Name:  "setting",
		Text:  "Setting",
		Value: "setting",
		Style: "default",
		Type:  "button",
	}

	listBtn := slack.AttachmentAction{
		Name:  "list-account",
		Text:  "List Account",
		Value: "list-account",
		Style: "primary",
		Type:  "button",
	}

	closeBtn := slack.AttachmentAction{
		Name:  "close",
		Text:  "Close",
		Value: "close",
		Style: "danger",
		Type:  "button",
	}

	at.Text = "Hi ! Happy to see you."
	at.CallbackID = "main-menu"
	at.Actions = []slack.AttachmentAction{addBtn, listBtn, settingBtn, closeBtn}

	return slack.Msg{
		Attachments: []slack.Attachment{at},
	}
}
