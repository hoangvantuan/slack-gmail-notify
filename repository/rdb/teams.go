package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Team is teams table
type Team struct {
	TeamID         string `gorm:"primary_key"`
	TeamName       string
	UserID         string
	Scope          string
	AccessToken    string
	BotUserID      string
	BotAccessToken string
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
	temp := &Team{}
	result := t.db.Where("team_id = ?", team.TeamID).First(temp)
	if result.RecordNotFound() {
		return t.db.Save(team).Error
	}
	if result.Error != nil {
		return result.Error
	}

	temp.AccessToken = team.AccessToken
	temp.BotAccessToken = team.BotAccessToken
	temp.BotUserID = team.BotUserID
	temp.Scope = team.Scope
	temp.UserID = team.UserID

	return t.db.Save(temp).Error
}
