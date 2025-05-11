package model

type Thumbnail struct {
	ID          string `json:"id"`
	VideoID     string `json:"video_id"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Format      string `json:"format"`
	Size        uint32 `json:"size"`
	TimeStamp   uint64 `json:"timestamp"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"`
	StoragePath string `json:"-"`
	IsDefault   bool   `json:"is_default"`
}
