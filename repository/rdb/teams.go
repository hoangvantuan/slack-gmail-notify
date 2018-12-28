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
	AccessToken    string `gorm:"not null;size:1000"`
	BotUserID      string
	BotAccessToken string `gorm:"size:1000"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// TeamRepository is defind interface for team
type TeamRepository interface {
	FindByTeamID(teamID string) (*Team, error)
	DeleteByTeamID(teamID string) error
	FindAllTeam() ([]*Team, error)
	Save(team *Team) error
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

	return team, nil
}

// DeleteByTeamID is delete team by id
func (t *teamRepositoryImpl) DeleteByTeamID(teamID string) error {
	return t.db.Where("team_id = ?", teamID).Delete(Team{}).Error
}

// FindAllTeam return all team info
func (t *teamRepositoryImpl) FindAllTeam() ([]*Team, error) {
	teams := []*Team{}

	result := t.db.Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}

	return teams, nil
}

// Save will save or update team
func (t *teamRepositoryImpl) Save(team *Team) error {
	return t.db.Save(team).Error
}
