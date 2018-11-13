package usecase

type eventUsecaseImpl struct{}

// EventUsecase is event interface
type EventUsecase interface {
	UninstallApp(teamID string) error
}

// NewEventUsecase will return event usecase
func NewEventUsecase() AuthUsecase {
	return &authUsecaseImpl{}
}

// UninstallApp will remove all data of team
func (e *eventUsecaseImpl) UninstallApp(teamID string) error {
	return nil
}
