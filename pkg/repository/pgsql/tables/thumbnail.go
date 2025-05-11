package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Thumbnail struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	VideoID     uuid.UUID      `gorm:"type:uuid;not null;index"`
	Width       int            `gorm:"not null"`
	Height      int            `gorm:"not null"`
	Format      string         `gorm:"not null"`
	Size        uint32         `gorm:"not null"` // Size in bytes
	TimeStamp   uint64         `gorm:"not null"` // Position in video where thumbnail was taken
	CreatedAt   time.Time      `gorm:"autoCreateTime:nano"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime:nano"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	StoragePath string         `gorm:"default:''"`
	IsDefault   bool           `gorm:"default:false;not null"`
}

func (Thumbnail) TableName() string {
	return "thumbnails"
}
