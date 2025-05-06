package tables

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Video struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Title           string         `gorm:"not null;unique" json:"title"`
	Description     string         `gorm:"not null" json:"description"`
	ParentID        *uuid.UUID     `gorm:"type:uuid" json:"parent_id,omitempty"` // Should be nullable for original videos
	Width           uint32         `json:"width"`                                // Will be unknown during upload
	Height          uint32         `json:"height"`                               // Will be unknown during upload
	UserID          uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Format          string         `json:"format"`            // Will be unknown initially
	Length          uint64         `json:"length"`            // Will be unknown during upload
	AudioSampleRate uint32         `json:"audio_sample_rate"` // Will be unknown during upload
	AudioCodec      string         `json:"audio_codec"`       // Will be unknown during upload
	RetryCount      uint8          `gorm:"default:0" json:"retry_count"`
	Status          string         `gorm:"not null" json:"status"`
	CreatedAt       time.Time      `gorm:"autoCreateTime:nano" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime:nano" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Should be nullable
	IsFeatured      bool           `gorm:"default:false" json:"is_featured,omitempty"`
	Visibility      string         `gorm:"not null" json:"visibility"`
	Slug            string         `gorm:"unique;not null" json:"slug"`
	Size            uint64         `json:"size"`     // Will be unknown during initial upload
	Language        string         `json:"language"` // Might be unknown initially
	StoragePath     string         `gorm:"default:''" json:"storage_path"`
}
