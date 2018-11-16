package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/nlopes/slack"
	"github.com/pkg/errors"
)

// IteractiveRequestParams is request from command
type IteractiveRequestParams struct {
	Type    string `json:"type"`
	Actions []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Value string `json:"value"`
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
	OpenSettingDialog(ir *IteractiveRequestParams) error
	ListAccount(ir *IteractiveRequestParams) error
}

// NewIteractiveUsecase will return event usecase
func NewIteractiveUsecase() IteractiveUsecase {
	return &iteractiveUsecaseImpl{}
}

func (i *iteractiveUsecaseImpl) OpenSettingDialog(ir *IteractiveRequestParams) error {
	slAPI, err := slackAPI(ir.Team.ID)
	if err != nil {
		return errors.Wrap(err, "have error while get slack client")
	}

	dl, err := settingDialog(ir)
	if err != nil {
		return errors.Wrap(err, "have error while generate dialog")
	}

	err = slAPI.OpenDialog(ir.TriggerID, *dl)
	if err != nil {
		infra.Swarn("has error while open dialog", err)
	}

	return nil
}

func (i *iteractiveUsecaseImpl) ListAccount(ir *IteractiveRequestParams) error {
	return nil
}

func listAccount(ir *IteractiveRequestParams) {

}

func settingDialog(ir *IteractiveRequestParams) (*slack.Dialog, error) {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	mails, err := gmailRepo.FindByUserID(ir.User.ID)
	if err != nil {
		// return empty dialog
		return nil, errors.Wrap(err, "have error while get list gmail")
	}

	elements := []slack.DialogElement{}
	mailsOption := []slack.DialogSelectOption{}
	for _, mail := range mails {
		mailsOption = append(mailsOption, slack.DialogSelectOption{
			Label: "Email",
			Value: mail.Email,
		})
	}

	element := slack.DialogInputSelect{
		DialogInput: slack.DialogInput{
			Type:  slack.InputTypeSelect,
			Label: "Email",
			Name:  "Email",
		},
		SelectedOptions: mails[0].Email,
		Options:         mailsOption,
	}

	elements = append(elements, element)

	return &slack.Dialog{
		CallbackID:  "setting-dialog",
		Title:       "Email list",
		SubmitLabel: "Change",
		Elements:    elements,
	}, nil
}
