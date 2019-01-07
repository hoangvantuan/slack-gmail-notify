package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
)

type eventUsecaseImpl struct{}

// EventUsecase is event interface
type EventUsecase interface {
	UninstallApp(teamID string) error
}

// NewEventUsecase will return event usecase
func NewEventUsecase() EventUsecase {
	return &eventUsecaseImpl{}
}

// UninstallApp will remove all data of team
func (e *eventUsecaseImpl) UninstallApp(teamID string) (err error) {
	rdbexRepository := rdb.NewRdbexRepository(infra.RDB)
	return rdbexRepository.DeleteTeam(teamID)
}
