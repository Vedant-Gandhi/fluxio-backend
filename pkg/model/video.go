package model

import (
	"time"

	"github.com/google/uuid"
)

type VideoStatus string

const (
	VideoStatusPending       VideoStatus = "pending"
	VideoStatusProcessing    VideoStatus = "processing"
	VideoGeneratingThumbnail VideoStatus = "generating_thumbnail"
	VideoStatusCompleted     VideoStatus = "completed"
	VideoStatusFailed        VideoStatus = "failed"
	VideoStatusDeleted       VideoStatus = "deleted"
	VideoStatusAbandoned     VideoStatus = "abandoned"
)

type Video struct {
	ID              uuid.UUID   `json:"id"`
	Title           string      `json:"title"`
	Description     string      `json:"description"`
	ParentID        *uuid.UUID  `json:"parent_id,omitempty"`
	Width           uint32      `json:"width"`
	Height          uint32      `json:"height"`
	UserID          uuid.UUID   `json:"user_id"`
	Format          string      `json:"format"`
	Length          uint64      `json:"length"`
	AudioSampleRate uint32      `json:"audio_sample_rate"`
	AudioCodec      string      `json:"audio_codec"`
	RetryCount      uint8       `json:"retry_count"`
	Status          VideoStatus `json:"status"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	DeletedAt       *time.Time  `json:"deleted_at,omitempty"`
	IsFeatured      bool        `json:"is_featured,omitempty"`
	Visibility      string      `json:"visibility"`
	Slug            string      `json:"slug"`
	Size            uint64      `json:"size"`
	Language        string      `json:"language"`
}
