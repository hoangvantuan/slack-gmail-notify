package rdb

import (
	"time"
)

// User is users table
type User struct {
	UserID    string `gorm:"primary_key"`
	TeamID    string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
