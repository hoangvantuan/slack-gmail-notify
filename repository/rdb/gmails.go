package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Gmail is gmails table
type Gmail struct {
	Email           string    `gorm:"primary_key"`
	UserID          string    `gorm:"not null"`
	AccessToken     string    `gorm:"not null;size:1000"`
	RefreshToken    string    `gorm:"not null;size:1000"`
	TokenType       string    `gorm:"not null"`
	Scope           string    `gorm:"not null"`
	ExpiryDate      time.Time `gorm:"not null"`
	NotifyChannelID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// GmailRepository is defind interface for team
type GmailRepository interface {
	DeleteByEmail(email string) error
	DeleteAllByUserID(userID string) error
	FindByEmail(email string) (*Gmail, error)
	FindByUserID(userID string) ([]*Gmail, error)
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
	result := t.db.First(gmail, email)
	if result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// DeleteByID is delete gmail by id
func (t *gmailRepositoryImpl) DeleteByEmail(email string) error {
	return t.db.Where("email = ?", email).Delete(Gmail{}).Error
}

// DeleteAllByUserID delete all user by user id
func (t *gmailRepositoryImpl) DeleteAllByUserID(userID string) error {
	result := t.db.Where("user_id = ?", userID).Delete(Gmail{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// FindByUserID will find gmail by id
func (t *gmailRepositoryImpl) FindByUserID(userID string) ([]*Gmail, error) {
	gmails := []*Gmail{}
	result := t.db.Where("user_id = ?", userID).Find(&gmails)

	if result.Error != nil {
		return nil, result.Error
	}

	return gmails, nil
}

// Save -
func (t *gmailRepositoryImpl) Save(gmail *Gmail) error {
	return t.db.Save(gmail).Error
}
