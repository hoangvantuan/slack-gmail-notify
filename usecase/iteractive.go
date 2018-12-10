package usecase

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/mdshun/slack-gmail-notify/worker"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// IteractiveRequestParams is request from command
type IteractiveRequestParams struct {
	Type    string `json:"type"`
	Actions []struct {
		Name            string                         `json:"name"`
		Type            string                         `json:"type"`
		Value           string                         `json:"value"`
		SelectedOptions []slack.AttachmentActionOption `json:"selected_options"`
	}
	CallbackID string `json:"callback_id"`
	Team       struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
	Channel struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"channel"`
	User struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
	ActionTs     string `json:"action_ts"`
	MessageTs    string `json:"message_ts"`
	AttachmentID string `json:"attachment_id"`
	Token        string `json:"token"`
	IsAppUnfurl  bool   `json:"is_app_unfurl"`
	ResponseURL  string `json:"response_url"`
	TriggerID    string `json:"trigger_id"`
}

type iteractiveUsecaseImpl struct{}

// IteractiveUsecase is event interface
type IteractiveUsecase interface {
	ListAccount(ir *IteractiveRequestParams) error
	NotifyChannel(ir *IteractiveRequestParams) error
	RemoveAccount(ir *IteractiveRequestParams) error
}

// NewIteractiveUsecase will return event usecase
func NewIteractiveUsecase() IteractiveUsecase {
	return &iteractiveUsecaseImpl{}
}

func (i *iteractiveUsecaseImpl) ListAccount(ir *IteractiveRequestParams) error {
	msg, err := listAccount(ir, "List gmail account you already registered")
	if err != nil {
		return errors.Wrap(err, "have error while get list acocunt")
	}

	msgjson, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "have error while marshal json")
	}

	_, err = http.Post(ir.ResponseURL, "application/json", bytes.NewReader(msgjson))
	if err != nil {
		return errors.Wrap(err, "have error while post message")
	}

	return nil
}

// TODO: need update worker
func (i *iteractiveUsecaseImpl) NotifyChannel(ir *IteractiveRequestParams) error {
	infra.Sdebug("notify account", ir)

	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	gmailID, err := strconv.Atoi(ir.CallbackID)
	if err != nil {
		return errors.Wrap(err, "can not convert gmail id")
	}

	gmail, err := gmailRepo.FindByID(gmailID)
	if err != nil {
		return errors.Wrap(err, "can not fetch gmail")
	}
	infra.Sdebug(gmail)

	gmail.NotifyChannelID = ir.Actions[0].SelectedOptions[0].Value

	_, err = gmailRepo.Update(gmail)
	if err != nil {
		return errors.Wrap(err, "can not update gmail")
	}

	slackAPI, err := util.SlackAPI(ir.Team.ID)
	if err != nil {
		return errors.Wrap(err, "error while init slack client")
	}

	err = worker.NotifyGmail(gmail, slackAPI)
	if err != nil {
		return errors.Wrap(err, "error while notify gmail for new channel")
	}

	return nil
}

// TODO: need remove worker
func (i *iteractiveUsecaseImpl) RemoveAccount(ir *IteractiveRequestParams) error {
	infra.Sdebug("remove account", ir)

	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	gmailID, err := strconv.Atoi(ir.CallbackID)
	if err != nil {
		return errors.Wrap(err, "can not convert gmail id")
	}

	err = gmailRepo.DeleteByID(gmailID)
	if err != nil {
		return errors.Wrap(err, "can not delete email with id")
	}

	// Stop notify gmail
	worker.StopNotifyGmail(&rdb.Gmail{
		ID: gmailID,
	})

	msg, err := listAccount(ir, "List gmail account you already register")
	if err != nil {
		return errors.Wrap(err, "have error while get list acocunt")
	}

	msgjson, err := json.Marshal(msg)
	if err != nil {
		return errors.Wrap(err, "have error while marshal json")
	}

	_, err = http.Post(ir.ResponseURL, "application/json", bytes.NewReader(msgjson))
	if err != nil {
		return errors.Wrap(err, "have error while post message")
	}

	return nil
}

func listAccount(ir *IteractiveRequestParams, text string) (*slack.Msg, error) {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	mails, err := gmailRepo.FindByUserID(ir.User.ID)
	if err != nil {
		// return empty dialog
		return nil, errors.Wrap(err, "have error while get list gmail")
	}

	if len(mails) == 0 {
		text = "You no have any gmail account, please add to start notify"
	}

	selectChannelBtn := func(value string) slack.AttachmentAction {
		return slack.AttachmentAction{
			Name:       "notify-channel",
			Text:       "Notify To",
			Type:       "select",
			DataSource: "channels",
			SelectedOptions: []slack.AttachmentActionOption{
				{
					Text:  value,
					Value: value,
				},
			},
		}
	}

	// value is gmail id
	removeBtn :=
		slack.AttachmentAction{
			Name:  "remove-gmail",
			Text:  "Remove",
			Style: "danger",
			Type:  "button",
		}

	closeBtn := slack.AttachmentAction{
		Name:  "close",
		Text:  "Close",
		Value: "close",
		Style: "danger",
		Type:  "button",
	}

	closeAt := slack.Attachment{
		CallbackID: "close",
		Actions: []slack.AttachmentAction{
			closeBtn,
		},
	}

	ats := []slack.Attachment{}

	for _, email := range mails {
		at := slack.Attachment{}
		at.Text = email.Email
		at.CallbackID = strconv.Itoa(email.ID)
		at.Actions = []slack.AttachmentAction{
			selectChannelBtn(email.NotifyChannelID),
			removeBtn,
		}

		infra.Sdebug(at)

		ats = append(ats, at)
	}

	ats = append(ats, closeAt)

	return &slack.Msg{
		Text:            text,
		ReplaceOriginal: true,
		Attachments:     ats,
	}, nil
}
