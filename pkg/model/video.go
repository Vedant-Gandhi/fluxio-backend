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
	VideoStatusUploadPending   VideoStatus = "upload_pending"
	VideoStatusProcessing      VideoStatus = "processing"
	VideoStatusProcessingDelay VideoStatus = "processing_delay"
	VideoStatusCompleted       VideoStatus = "completed"
	VideoStatusFailed          VideoStatus = "failed"
	VideoStatusDeleted         VideoStatus = "deleted"
	VideoStatusAbandoned       VideoStatus = "abandoned"

	VideoVisibilityPublic  VideoVisibility = "public"
	VideoVisibilityPrivate VideoVisibility = "private"
)

// This function checks if the video status is of a valid value.
func (s VideoStatus) IsAcceptable() bool {
	switch s {
	case VideoStatusUploadPending,
		VideoStatusProcessingDelay,
		VideoStatusProcessing,
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

const (
	VidInternalStatusUploadPending       VideoInternalStatus = "upload_pending"
	VidInternalStatusMetaExtracted       VideoInternalStatus = "meta_extracted"
	VidInternalStatusThumbnailGenerated  VideoInternalStatus = "thumbnail_generated"
	VidInternalStatusProcessingCompleted VideoInternalStatus = "completed"

	VidInternalStatusThumbnailFailed VideoInternalStatus = "thumbnail_failed"
	VidInternalStatusMetaFailed      VideoInternalStatus = "meta_failed"
)

type VideoInternalStatus string

// This function checks if the video status is of a valid value.
func (s VideoInternalStatus) IsAcceptable() bool {
	switch s {
	case VidInternalStatusUploadPending:
		return true
	default:
		return false
	}
}

func (s VideoInternalStatus) String() string {
	return string(s)
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
	CreatedAt       *time.Time      `json:"created_at"`
	UpdatedAt       *time.Time      `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
	IsFeatured      bool            `json:"is_featured,omitempty"`
	Visibility      VideoVisibility `json:"visibility"`
	Slug            string          `json:"slug"`
	Size            float32         `json:"size"`
	Language        string          `json:"language"`
	ResourceURL     url.URL         `json:"resource_url"`
	StoragePath     string          `json:"-"`
	Thumbnails      []Thumbnail     `json:"thumbnails,omitempty"`
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
	InternalStatus  VideoStatus     `json:"-"`
	IsFeatured      bool            `json:"is_featured,omitempty"`
	Visibility      VideoVisibility `json:"visibility"`
	Slug            string          `json:"slug"`
	Size            float32         `json:"size"`
	Language        string          `json:"language"`
	ResourceURL     url.URL         `json:"resource_url"`
	StoragePath     string          `json:"-"`
}
