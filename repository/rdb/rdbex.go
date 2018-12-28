package rdb

import (
	"github.com/jinzhu/gorm"
)

type RdbexRepository interface {
	DeleteTeam(teamId string) error
}

type rdbexRepositoryImpl struct {
	db *gorm.DB
}

func NewRdbexRepository(db *gorm.DB) RdbexRepository {
	return &rdbexRepositoryImpl{db}
}

func (r *rdbexRepositoryImpl) DeleteTeam(teamID string) error {
	sql = `
	DELETE teams as t, users as u, gmails as g
	 
	`
}
