package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/mdshun/slack-gmail-notify/util"
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
}

// NewIteractiveUsecase will return event usecase
func NewIteractiveUsecase() IteractiveUsecase {
	return &iteractiveUsecaseImpl{}
}

func (i *iteractiveUsecaseImpl) OpenSettingDialog(ir *IteractiveRequestParams) error {
	teamRepo := rdb.NewTeamRepository(infra.RDB)

	team, err := teamRepo.FindByTeamID(ir.Team.ID)
	if err != nil {
		return errors.Wrap(err, errWhileFindTeam)
	}

	token, _ := util.Decrypt(team.BotAccessToken, infra.Env.EncryptKey)
	slackAPI := slack.New(token)

	err = slackAPI.OpenDialog(ir.TriggerID, settingDialog(slackAPI, ir))
	if err != nil {
		infra.Swarn("has error while open dialog", err)
	}

	return nil
}

func settingDialog(api *slack.Client, ir *IteractiveRequestParams) slack.Dialog {
	gmailRepo := rdb.NewGmailRepository(infra.RDB)
	mails, err := gmailRepo.FindByUserID(ir.User.ID)
	if err != nil {
		// return empty dialog
		return slack.Dialog{}
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

	return slack.Dialog{
		CallbackID:  "setting-dialog",
		Title:       "Email list",
		SubmitLabel: "Change",
		Elements:    elements,
	}
}
