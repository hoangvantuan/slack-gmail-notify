package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository"
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
	tx := infra.RDB.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	teamRepo := repository.NewTeamRepository(tx)
	userRepo := repository.NewUserRepository(tx)
	gmailRepo := repository.NewGmailRepository(tx)

	err = teamRepo.DeleteByTeamID(teamID)

	if err != nil {
		infra.Swarn(errWhileDeleteTeam, err)
		return
	}

	users, err := userRepo.FindAllByTeamID(teamID)

	if err != nil {
		infra.Swarn(errWhileFindUser, err)
		return
	}

	for _, user := range users {
		err = gmailRepo.DeleteAllByUserID(user.UserID)

		if err != nil {
			infra.Swarn(errWhileDeleteGmail, err)
			return
		}
	}

	err = userRepo.DeleteAllByTeamID(teamID)

	if err != nil {
		infra.Swarn(errWhileDeleteUser, err)
		return
	}

	tx.Commit()

	return nil
}
