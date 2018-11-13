package usecase

// AuthRequestInput is auth request param
type EventRequestInput struct {
	Code  string
	State string
}

type eventUsecaseImpl struct{}

// AuthUsecase is auth interface
type EventUsecase interface {
	SlackAuth(ri *AuthRequestInput) error
}

// NewAuthUsecase will return auth usecase
func NewEventUsecase() AuthUsecase {
	return &authUsecaseImpl{}
}