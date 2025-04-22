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

// Repo Errors
var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
	ErrUserNotFound      = errors.New("user not found")
)
