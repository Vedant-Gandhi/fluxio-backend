package tables

import (
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Email     string    `gorm:"unique;not null" json:"email"`
	CreatedAt string    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt string    `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt string    `gorm:"autoDeleteTime" json:"deleted_at"`
}
