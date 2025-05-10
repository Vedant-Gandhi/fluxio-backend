package constants

import "time"

// Auth related constants
const (
	AuthTokenCookieName = "auth_token"
	AuthTokenCookieExp  = int(8 * time.Hour) // 8 hours
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
