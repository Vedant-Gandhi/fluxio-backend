package constants

import "time"

// Auth related constants
const (
	AuthTokenCookieName = "auth_token"
	AuthTokenCookieExp  = int(8 * 60 * 60) // 8 hours
)

// Gin related constants
const (
	GinUserContextKey = "user"
)

const (
	MaxVideoURLRegenerateRetryCount       = 4
	MaxVideoThumbnailRegenerateRetryCount = 3
)

const (
	PreSignedVidUploadURLExpireTime       = 1 * time.Hour
	PreSignedVidTempDownloadURLExpireTime = 1 * time.Hour
)

const (
	VidSizeDecimalPrecision = 3
	TotalThumbnailCount     = 3
)

var (

	// Declare it here to reduce creation costs for the GC.
	ValidVideoMimes = []string{
		"video/mp4",        // MP4 format
		"video/webm",       // WebM format
		"video/ogg",        // Ogg video
		"video/3gpp",       // 3GPP mobile video
		"video/quicktime",  // MOV format
		"video/x-msvideo",  // AVI format
		"video/mpeg",       // MPEG format
		"video/x-matroska", // MKV format
		"video/mp2t",       // MPEG Transport Stream (.ts files)
	}
)
