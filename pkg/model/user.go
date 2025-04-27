package model

import (
	"time"
)

type UserID string

func (id UserID) String() string {
	return string(id)
}

type User struct {
	ID            UserID    `json:"id,omitempty"`
	Username      string    `json:"username"`
	Password      string    `json:"-"` // Hide password in JSON responses
	Email         string    `json:"email"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	IsBlackListed bool      `json:"is_blacklisted"`
}
