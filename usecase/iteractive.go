package usecase

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
}

// NewIteractiveUsecase will return event usecase
func NewIteractiveUsecase() IteractiveUsecase {
	return &iteractiveUsecaseImpl{}
}
