package model

import "time"

type ThumbnailID string

func (t ThumbnailID) String() string {
	return string(t)
}

type Thumbnail struct {
	ID          ThumbnailID `json:"id"`
	VideoID     VideoID     `json:"video_id"`
	Width       uint16      `json:"width"`
	Height      uint16      `json:"height"`
	Format      string      `json:"format"`
	Size        uint32      `json:"size"`
	TimeStamp   uint64      `json:"timestamp"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	DeletedAt   time.Time   `json:"deleted_at,omitempty"`
	StoragePath string      `json:"-"`
	IsDefault   bool        `json:"is_default"`
}
