package models

import (
	"time"
)

// User is the application user's details. This should never be used outside the
// application, use Profile instead.
type User struct {
	ID        int64
	Name      string `gorm:"index:idx_username"`
	Password  string
	Profie    Profile
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Profile is the details about the application user.
type Profile struct {
	ID        int64
	UserID    int64
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
