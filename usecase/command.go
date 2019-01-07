package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/nlopes/slack"
)

type commandUsecaseImpl struct{}

// CommandUsecase is event interface
type CommandUsecase interface {
	GetMainMenu(rp *slack.SlashCommand) error
}

// NewCommandUsecase will return event usecase
func NewCommandUsecase() CommandUsecase {
	return &commandUsecaseImpl{}
}

func (c *commandUsecaseImpl) GetMainMenu(rp *slack.SlashCommand) error {
	msgAt, err := genInteractiveMenu(rp, "Hi, can i help you ?")
	if err != nil {
		return err
	}

	msgatstr, err := json.Marshal(msgAt)
	if err != nil {
		return err
	}

	_, err = http.Post(rp.ResponseURL, "application/json", bytes.NewReader(msgatstr))
	if err != nil {
		return err
	}

	return nil
}

func genInteractiveMenu(rp *slack.SlashCommand, text string) (*slack.Msg, error) {
	at := slack.Attachment{}

	pjson, err := json.Marshal(rp)
	if err != nil {
		return nil, err
	}
	pjsoneconded, err := util.Encrypt(string(pjson))
	if err != nil {
		return nil, err
	}
	// generate action
	addBtn := slack.AttachmentAction{
		Name:  util.AddGmailAccountName,
		Text:  util.AddGmailAccountText,
		Value: util.AddGmailAccountValue,
		Style: util.AddGmailAccountStyle,
		Type:  util.AddGmailAccountType,
		URL:   fmt.Sprintf("%s/v1/auth/google?state=%s", infra.Env.APIHost, pjsoneconded),
	}

	listBtn := slack.AttachmentAction{
		Name:  util.ListGmailAccountName,
		Text:  util.ListGmailAccountText,
		Value: util.ListGmailAccountValue,
		Style: util.ListGmailAccountStyle,
		Type:  util.ListGmailAccountType,
	}

	closeBtn := slack.AttachmentAction{
		Name:  util.CloseName,
		Text:  util.CloseText,
		Value: util.CloseValue,
		Style: util.CloseStyle,
		Type:  util.CloseType,
	}

	at.CallbackID = "main-menu"
	at.Actions = []slack.AttachmentAction{addBtn, listBtn, closeBtn}

	return &slack.Msg{
		Text:            text,
		ReplaceOriginal: true,
		Attachments:     []slack.Attachment{at},
	}, nil
}
