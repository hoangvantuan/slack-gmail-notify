package repository

import (
	"database/sql"
	"time"

	"github.com/jinzhu/gorm"
)

// User is users table
type User struct {
	ID        uint           `gorm:"primary_key"`
	UserID    sql.NullString `gorm:"not null"`
	TeamID    sql.NullString `gorm:"not null"`
	UserName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository is defind interface for team
type UserRepository interface {
	FindByID(id uint) (*User, error)
	Add(user *User) (*User, error)
	DeleteByID(id uint) error
	Update(user *User) (*User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

// NewUserRepository is return data access object for user table
func NewUserRepository(db *gorm.DB) UserRepository {
	// Migrate table if not exist
	if !db.HasTable(&User{}) {
		db.AutoMigrate(&User{})
	}

	return &userRepositoryImpl{db}
}

// FindByID will find user by id
func (t *userRepositoryImpl) FindByID(id uint) (*User, error) {
	user := &User{}
	result := t.db.First(user, id)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RecordNotFound() {
		return nil, ErrRecordNotFound
	}

	return user, nil
}

// Add is add user record
func (t *userRepositoryImpl) Add(user *User) (*User, error) {
	result := t.db.Create(user)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.NewRecord(user) {
		return nil, ErrCanNotCreateRecord
	}

	return user, nil
}

// DeleteByID is delete user by id
func (t *userRepositoryImpl) DeleteByID(id uint) error {
	result := t.db.Where("id = ?", id).Delete(User{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Update user
func (t *userRepositoryImpl) Update(user *User) (*User, error) {
	result := t.db.Update(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
