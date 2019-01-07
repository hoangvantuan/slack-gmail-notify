package usecase

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
	"github.com/mdshun/slack-gmail-notify/worker"
	"github.com/nlopes/slack"
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
	ListAllAccount(ir *IteractiveRequestParams) error
	NotifyToChannel(ir *IteractiveRequestParams) error
	RemoveAccount(ir *IteractiveRequestParams) error
}

// NewIteractiveUsecase will return event usecase
func NewIteractiveUsecase() IteractiveUsecase {
	return &iteractiveUsecaseImpl{}
}

func (i *iteractiveUsecaseImpl) ListAllAccount(ir *IteractiveRequestParams) error {
	msg, err := listAccount(ir, "List gmail account you already registered")
	if err != nil {
		return err
	}

	msgjson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = http.Post(ir.ResponseURL, "application/json", bytes.NewReader(msgjson))
	if err != nil {
		return err
	}

	return nil
}

func (i *iteractiveUsecaseImpl) NotifyToChannel(ir *IteractiveRequestParams) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	// CallbackID is email
	gmail, err := gmailRepo.FindByEmail(ir.CallbackID)
	if err != nil {
		return err
	}

	gmail.NotifyChannelID = ir.Actions[0].SelectedOptions[0].Value

	err = gmailRepo.Save(gmail)
	if err != nil {
		return err
	}

	teamRepo := rdb.NewTeamRepository(infra.RDB)
	team, err := teamRepo.FindByTeamID(ir.Team.ID)
	if err != nil {
		return err
	}

	slackAPI := slack.New(team.BotAccessToken)

	err = worker.NotifyForGmail(gmail, slackAPI)
	if err != nil {
		return err
	}

	return nil
}

func (i *iteractiveUsecaseImpl) RemoveAccount(ir *IteractiveRequestParams) error {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)

	err := gmailRepo.DeleteByEmail(ir.CallbackID)
	if err != nil {
		return err
	}

	// Stop notify gmail
	worker.StopNotifyForGmail(ir.CallbackID)

	msg, err := listAccount(ir, "List gmail account you already register")
	if err != nil {
		return err
	}

	msgjson, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = http.Post(ir.ResponseURL, "application/json", bytes.NewReader(msgjson))
	if err != nil {
		return err
	}

	return nil
}

func listAccount(ir *IteractiveRequestParams, text string) (*slack.Msg, error) {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	mails, err := gmailRepo.FindByUser(&rdb.User{
		TeamID: ir.Team.ID,
		UserID: ir.User.ID,
	})

	if err != nil {
		// return empty dialog
		return nil, err
	}

	if len(mails) == 0 {
		text = "You not have any gmail account, please add to start notify"
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

	closeAt := slack.Attachment{
		CallbackID: "close",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  util.CloseName,
				Text:  util.CloseText,
				Value: util.CloseValue,
				Style: util.CloseStyle,
				Type:  util.CloseType,
			},
		},
	}

	ats := []slack.Attachment{}

	for _, email := range mails {
		at := slack.Attachment{}
		at.Text = email.Email
		// callback_id is email
		at.CallbackID = email.Email
		at.Actions = []slack.AttachmentAction{
			selectChannelBtn(email.NotifyChannelID),
			removeBtn,
		}

		ats = append(ats, at)
	}

	ats = append(ats, closeAt)

	return &slack.Msg{
		Text:            text,
		ReplaceOriginal: true,
		Attachments:     ats,
	}, nil
}
