package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Gmail is gmails table
type Gmail struct {
	ID              int       `gorm:"primary_key"`
	UserID          string    `gorm:"not null"`
	Email           string    `gorm:"not null"`
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
	FindByID(id int) (*Gmail, error)
	Add(gmail *Gmail) (*Gmail, error)
	DeleteByID(id int) error
	Update(gmail *Gmail) (*Gmail, error)
	DeleteAllByUserID(userID string) error
	FindByEmail(email, userID string) (*Gmail, error)
	FindByUserID(userID string) ([]Gmail, error)
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
func (t *gmailRepositoryImpl) FindByID(id int) (*Gmail, error) {
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
func (t *gmailRepositoryImpl) DeleteByID(id int) error {
	result := t.db.Where("id = ?", id).Delete(Gmail{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update gmail
func (t *gmailRepositoryImpl) Update(gmail *Gmail) (*Gmail, error) {
	result := t.db.Save(gmail)

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

// FindByEmail will find gmail by id
func (t *gmailRepositoryImpl) FindByEmail(email, userID string) (*Gmail, error) {
	gmail := &Gmail{}
	result := t.db.Where("email = ? AND user_id = ?", email, userID).First(gmail)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}

	return gmail, nil
}

// FindByUserID will find gmail by id
func (t *gmailRepositoryImpl) FindByUserID(userID string) ([]Gmail, error) {
	gmails := []Gmail{}
	result := t.db.Where("user_id = ?", userID).Find(&gmails)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}

	return gmails, nil
}
