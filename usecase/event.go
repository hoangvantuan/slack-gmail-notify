package usecase

import (
	"github.com/mdshun/slack-gmail-notify/infra"
	"github.com/mdshun/slack-gmail-notify/repository/rdb"
	"github.com/pkg/errors"
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
		} else {
			tx.Commit()
		}
	}()

	teamRepo := rdb.NewTeamRepository(tx)
	userRepo := rdb.NewUserRepository(tx)
	gmailRepo := rdb.NewGmailRepository(tx)

	err = teamRepo.DeleteByTeamID(teamID)
	if err != nil {
		return errors.Wrap(err, "have error while delete team")
	}

	users, err := userRepo.FindAllByTeamID(teamID)
	if err != nil {
		return errors.Wrap(err, "have error while find user")
	}

	for _, user := range users {
		err = gmailRepo.DeleteAllByUserID(user.UserID)
		if err != nil {
			return errors.Wrap(err, "have error while delete email")
		}
	}

	err = userRepo.DeleteAllByTeamID(teamID)
	if err != nil {
		return errors.Wrap(err, "have error while delete user")
	}

	return nil
}
