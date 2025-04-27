package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username      string         `gorm:"unique;not null" json:"username"`
	Password      string         `gorm:"not null" json:"-"` // Hide password in JSON responses
	Email         string         `gorm:"unique;not null" json:"email"`
	CreatedAt     time.Time      `gorm:"autoCreateTime:nano" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	IsBlackListed bool           `gorm:"default:false" json:"is_blacklisted,omitempty"`
}
