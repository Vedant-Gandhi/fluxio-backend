package model

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type VideoID string

type VideoStatus string

type VideoVisibility string

const (
	VideoStatusPending       VideoStatus = "pending"
	VideoStatusProcessing    VideoStatus = "processing"
	VideoStatusMetaExtracted VideoStatus = "meta_extracted"
	VideoGeneratingThumbnail VideoStatus = "generating_thumbnail"
	VideoStatusCompleted     VideoStatus = "completed"
	VideoStatusFailed        VideoStatus = "failed"
	VideoStatusDeleted       VideoStatus = "deleted"
	VideoStatusAbandoned     VideoStatus = "abandoned"

	VideoVisibilityPublic  VideoVisibility = "public"
	VideoVisibilityPrivate VideoVisibility = "private"
)

// This function checks if the video status is of a valid value.
func (s VideoStatus) IsAcceptable() bool {
	switch s {
	case VideoStatusPending,
		VideoStatusProcessing,
		VideoStatusMetaExtracted,
		VideoGeneratingThumbnail,
		VideoStatusCompleted,
		VideoStatusFailed,
		VideoStatusDeleted,
		VideoStatusAbandoned:
		return true
	default:
		return false
	}
}

func (s VideoStatus) String() string {
	return string(s)
}

func (s VideoVisibility) String() string {
	return string(s)
}

func (s VideoID) String() string {
	return string(s)
}

// This function checks if the video visibility is of a valid value.
func (s VideoVisibility) IsAcceptable() bool {
	switch s {
	case VideoVisibilityPublic,
		VideoVisibilityPrivate:
		return true
	default:
		return false
	}
}

type Video struct {
	ID              VideoID         `json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	ParentID        *uuid.UUID      `json:"parent_id,omitempty"`
	Width           uint32          `json:"width"`
	Height          uint32          `json:"height"`
	UserID          uuid.UUID       `json:"user_id"`
	Format          string          `json:"format"`
	Length          uint64          `json:"length"`
	AudioSampleRate uint32          `json:"audio_sample_rate"`
	AudioCodec      string          `json:"audio_codec"`
	RetryCount      uint8           `json:"retry_count"`
	Status          VideoStatus     `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
	IsFeatured      bool            `json:"is_featured,omitempty"`
	Visibility      VideoVisibility `json:"visibility"`
	Slug            string          `json:"slug"`
	Size            float32         `json:"size"`
	Language        string          `json:"language"`
	ResourceURL     url.URL         `json:"resource_url"`
	StoragePath     string          `json:"-"`
	Thumnbails      []Thumbnail     `json:"thumbnails,omitempty"`
}

type UpdateVideoMeta struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	ParentID        *uuid.UUID      `json:"parent_id,omitempty"`
	Width           uint32          `json:"width"`
	Height          uint32          `json:"height"`
	Format          string          `json:"format"`
	Length          uint64          `json:"length"`
	AudioSampleRate uint32          `json:"audio_sample_rate"`
	AudioCodec      string          `json:"audio_codec"`
	RetryCount      uint8           `json:"retry_count"`
	Status          VideoStatus     `json:"status"`
	IsFeatured      bool            `json:"is_featured,omitempty"`
	Visibility      VideoVisibility `json:"visibility"`
	Slug            string          `json:"slug"`
	Size            float32         `json:"size"`
	Language        string          `json:"language"`
	ResourceURL     url.URL         `json:"resource_url"`
	StoragePath     string          `json:"-"`
}
