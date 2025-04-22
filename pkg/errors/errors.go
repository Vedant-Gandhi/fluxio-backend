package fluxerrors

import "errors"

var (
	ErrPasswordTooShort = errors.New("password too short (minimum 8 characters)")
)

var (
	ErrPasswordFailed = errors.New("password failed to hash")
)

var (
	ErrUserCreationFailed = errors.New("failed to create new user")
)

// Repo Errors

var (
	ErrUsernameExists    = errors.New("username already exists")
	ErrEmailExists       = errors.New("email already exists")
	ErrUserAlreadyExists = errors.New("user already exists")
)
