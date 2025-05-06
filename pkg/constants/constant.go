package constants

// Auth related constants
const (
	AuthTokenCookieName = "auth_token"
	AuthTokenCookieExp  = 8 * 3600 // 8 hours
)

// Gin related constants
const (
	GinUserContextKey = "user"
)

const (
	MaxVideoURLRegenerateRetryCount       = 4
	MaxVideoThumbnailRegenerateRetryCount = 3
)
