package rdb

import (
	"github.com/jinzhu/gorm"
)

// RdbexRepository -
type RdbexRepository interface {
	DeleteTeam(teamID string) error
}

type rdbexRepositoryImpl struct {
	db *gorm.DB
}

// NewRdbexRepository -
func NewRdbexRepository(db *gorm.DB) RdbexRepository {
	return &rdbexRepositoryImpl{db}
}

// DeleteTeam delete all team
func (r *rdbexRepositoryImpl) DeleteTeam(teamID string) error {
	sql := `
	DELETE teams, gmails
	FROM teams as t
	INNER JOIN gmails as g ON g.team_id = g.team_id
	WHERE t.team_id = ?
	`

	return r.db.Exec(sql, teamID).Error
}
