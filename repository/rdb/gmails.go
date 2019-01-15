package rdb

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// Gmail is gmails table
const (
	Read    = "read"
	Unread  = "unread"
	Stop    = "stop"
	Working = "working"
	Pending = "pending"
)

type Gmail struct {
	Email           string `gorm:"primary_key"`
	TeamID          string
	UserID          string
	AccessToken     string
	RefreshToken    string
	TokenType       string
	Scope           string
	ExpiryDate      time.Time
	NotifyChannelID string
	MarkAs          string
	Status          string
	LabelID         string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// GmailRepository is defind interface for team
type GmailRepository interface {
	DeleteByEmail(email string) error
	DeleteByUser(user *User) error
	FindByEmail(email string) (*Gmail, error)
	FindByTeamID(teamID string) ([]*Gmail, error)
	FindByUser(user *User) ([]*Gmail, error)
	Save(gmail *Gmail) error
}

type gmailRepositoryImpl struct {
	db *gorm.DB
}

// NewGmailRepository is return data access object for gmail table
func NewGmailRepository(db *gorm.DB) GmailRepository {
	// Migrate table if not exist
	if !db.HasTable(&Gmail{}) {
		db.AutoMigrate(&Gmail{})
	}

	return &gmailRepositoryImpl{db}
}

// FindByEmail will find gmail by id
func (t *gmailRepositoryImpl) FindByEmail(email string) (*Gmail, error) {
	gmail := &Gmail{}
	result := t.db.Where("email = ?", email).First(gmail)
	if result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// DeleteByID is delete gmail by id
func (t *gmailRepositoryImpl) DeleteByEmail(email string) error {
	return t.db.Where("email = ?", email).Delete(Gmail{}).Error
}

// DeleteByUser delete all user by user id
func (t *gmailRepositoryImpl) DeleteByUser(user *User) error {
	return t.db.Where("user_id = ? AND team_id = ?", user.UserID, user.TeamID).Delete(Gmail{}).Error
}

// FindByTeamID will find gmail by id
func (t *gmailRepositoryImpl) FindByTeamID(teamID string) ([]*Gmail, error) {
	gmails := []*Gmail{}
	result := t.db.Where("team_id = ?", teamID).Find(&gmails)

	if result.Error != nil {
		return nil, result.Error
	}

	return gmails, nil
}

// FindByUserIDAndTeamID -
func (t *gmailRepositoryImpl) FindByUser(user *User) ([]*Gmail, error) {
	gmails := []*Gmail{}
	result := t.db.Where("team_id = ? AND user_id = ?", user.TeamID, user.UserID).Find(&gmails)

	if result.Error != nil {
		return nil, result.Error
	}

	return gmails, nil
}

// Save -
func (t *gmailRepositoryImpl) Save(gmail *Gmail) error {
	temp := &Gmail{}
	result := t.db.Where("email = ?", gmail.Email).First(temp)
	if result.RecordNotFound() {
		return t.db.Save(gmail).Error
	}
	if result.Error != nil {
		return result.Error
	}

	if temp.TeamID != gmail.TeamID {
		return errors.New("email was install in other workspace before")
	}

	if temp.UserID != gmail.UserID {
		return errors.New("email was install by other user before")
	}

	temp.AccessToken = gmail.AccessToken
	temp.ExpiryDate = gmail.ExpiryDate
	temp.RefreshToken = gmail.RefreshToken
	temp.Scope = gmail.Scope
	temp.TokenType = gmail.TokenType
	temp.NotifyChannelID = gmail.NotifyChannelID
	temp.MarkAs = gmail.MarkAs
	temp.LabelID = gmail.LabelID
	temp.Status = gmail.Status

	return t.db.Save(temp).Error
}
