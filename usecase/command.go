package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository"
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
	teamRepo := repository.NewTeamRepository(infra.RDB)

	team, err := teamRepo.FindByTeamID(rp.TeamID)
	if err != nil {
		return errors.Wrap(err, errWhileFindTeam)
	}

	slackAPI := slack.New(team.BotAccessToken)

	msgAt := genInteractiveMenu(rp)

	ts, err := slackAPI.PostEphemeral(rp.ChannelID, rp.UserID, slack.MsgOptionAttachments(msgAt))

	if err != nil {
		infra.Swarn(errWhilePostMsg, err)
		return errors.Wrap(err, errWhilePostMsg)
	}

	infra.Sdebug("Message successfully sent to channel ", rp.ChannelID, " for user ", rp.UserID, " at ", ts)

	return nil
}

func genInteractiveMenu(rp *CommandRequestParams) slack.Attachment {
	at := slack.Attachment{}

	// generate action
	addBtn := slack.AttachmentAction{
		Name:  "AddGmailAccount",
		Text:  "Add Account",
		Value: "AddGmailAccount",
		Style: "primary",
		Type:  "button",
		URL:   infra.Env.APIHost + "/v1/auth/google",
	}

	settingBtn := slack.AttachmentAction{
		Name:  "setting",
		Text:  "Setting",
		Value: "setting",
		Style: "default",
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
	at.CallbackID = "MainMenu"
	at.Actions = []slack.AttachmentAction{addBtn, settingBtn, closeBtn}

	return at
}
