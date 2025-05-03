package fluxerrors

import "errors"

var (
	ErrPasswordTooShort = errors.New("password too short (minimum 8 characters)")
	ErrPasswordFailed   = errors.New("password failed to hash")
)

var (
	ErrUserCreationFailed = errors.New("failed to create new user")
	ErrUserLoginFailed    = errors.New("failed to login user")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// Generic errors
var (
	ErrUnknown = errors.New("unknown error")
)

// Repo Errors
var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrUserNotFound      = errors.New("user not found")
)

// JWT errors
var (
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenNotFound       = errors.New("token not found")
	ErrFailedToCreateToken = errors.New("failed to create token")
	ErrInvalidClaims       = errors.New("invalid claims in token")
)

// Video errors
var (
	ErrVideoNotFound       = errors.New("video not found")
	ErrVideoAlreadyExists  = errors.New("video already exists")
	ErrVideoCreationFailed = errors.New("failed to create video")

	ErrDuplicateVideoTitle = errors.New("video title already exists")

	ErrFailedToGenerateVideoSlug = errors.New("failed to generate video slug")

	ErrVideoURLGenerationFailed = errors.New("failed to generate video upload URL")
	ErrVideoUploadNotAllowed    = errors.New("video upload not allowed")
)
