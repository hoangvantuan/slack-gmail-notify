package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Team is teams table
type Team struct {
	TeamID         string `gorm:"primary_key"`
	TeamName       string `gorm:"not null"`
	UserID         string `gorm:"not null"`
	Scope          string `gorm:"not null"`
	AccessToken    string `gorm:"not null"`
	BotUserID      string
	BotAccessToken string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TeamRepository is defind interface for team
type TeamRepository interface {
	FindByTeamID(teamID string) (*Team, error)
	Add(team *Team) (*Team, error)
	DeleteByTeamID(teamID string) error
	Update(team *Team) (*Team, error)
}

type teamRepositoryImpl struct {
	db *gorm.DB
}

// NewTeamRepository is return data access object for team table
func NewTeamRepository(db *gorm.DB) TeamRepository {
	// Migrate table if not exist
	if !db.HasTable(&Team{}) {
		db.AutoMigrate(&Team{})
	}

	return &teamRepositoryImpl{db}
}

// FindByTeamID will find team by id
func (t *teamRepositoryImpl) FindByTeamID(teamID string) (*Team, error) {
	team := &Team{}
	result := t.db.Where("team_id = ?", teamID).First(team)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}

	return team, nil
}

// Add is add team record
func (t *teamRepositoryImpl) Add(team *Team) (*Team, error) {
	result := t.db.Create(team)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.NewRecord(team) {
		return nil, ErrCanNotCreateRecord
	}

	return team, nil
}

// DeleteByTeamID is delete team by id
func (t *teamRepositoryImpl) DeleteByTeamID(teamID string) error {
	result := t.db.Where("team_id = ?", teamID).Delete(Team{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update team
func (t *teamRepositoryImpl) Update(team *Team) (*Team, error) {
	result := t.db.Save(team)

	if result.Error != nil {
		return nil, result.Error
	}

	return team, nil
}
