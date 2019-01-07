package rdb

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User is users table
type User struct {
	UserID    string `gorm:"primary_key"`
	TeamID    string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// UserRepository is defind interface for team
type UserRepository interface {
	FindByID(uid, tid string) (*User, error)
	FindByTeamID(teamID string) ([]User, error)
	DeleteByTeamID(teamID string) error
	Save(user *User) error
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
func (t *userRepositoryImpl) FindByID(uid, tid string) (*User, error) {
	user := &User{}
	result := t.db.Where("user_id = ? AND team_id = ?", uid, tid).First(user)

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

// FindByTeamID return all user with team id
func (t *userRepositoryImpl) FindByTeamID(teamID string) ([]User, error) {
	users := []User{}

	result := t.db.Where("team_id = ?", teamID).Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

// DeleteByTeamID delete all user by teamid
func (t *userRepositoryImpl) DeleteByTeamID(teamID string) error {
	result := t.db.Where("team_id = ?", teamID).Delete(User{})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Save user
func (t *userRepositoryImpl) Save(user *User) error {
	return t.db.Save(user).Error
}
