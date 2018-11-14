package repository

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Gmail is gmails table
type Gmail struct {
	ID              uint   `gorm:"primary_key"`
	UserID          string `gorm:"not null"`
	Email           string `gorm:"unique_index;not null"`
	AccessToken     string `gorm:"not null"`
	RefreshToken    string `gorm:"not null"`
	TokenType       string `gorm:"not null"`
	Scope           string `gorm:"not null"`
	ExpiryDate      string `gorm:"not null"`
	NotifyChannelID string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// GmailRepository is defind interface for team
type GmailRepository interface {
	FindByID(id uint) (*Gmail, error)
	Add(gmail *Gmail) (*Gmail, error)
	DeleteByID(id uint) error
	Update(gmail *Gmail) (*Gmail, error)
	DeleteAllByUserID(userID string) error
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

// FindByID will find gmail by id
func (t *gmailRepositoryImpl) FindByID(id uint) (*Gmail, error) {
	gmail := &Gmail{}
	result := t.db.First(gmail, id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}

	return gmail, nil
}

// Add is add gmail record
func (t *gmailRepositoryImpl) Add(gmail *Gmail) (*Gmail, error) {
	result := t.db.Create(gmail)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.NewRecord(gmail) {
		return nil, ErrCanNotCreateRecord
	}

	return gmail, nil
}

// DeleteByID is delete gmail by id
func (t *gmailRepositoryImpl) DeleteByID(id uint) error {
	result := t.db.Where("id = ?", id).Delete(Gmail{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update gmail
func (t *gmailRepositoryImpl) Update(gmail *Gmail) (*Gmail, error) {
	result := t.db.Update(gmail)

	if result.Error != nil {
		return nil, result.Error
	}

	return gmail, nil
}

// DeleteAllByUserID delete all user by user id
func (t *gmailRepositoryImpl) DeleteAllByUserID(userID string) error {
	result := t.db.Where("user_id = ?", userID).Delete(Gmail{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
