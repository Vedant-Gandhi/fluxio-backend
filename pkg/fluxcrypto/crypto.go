package fluxcrypto

import (
	"crypto/sha256"
	"encoding/base64"
	fluxerrors "fluxio-backend/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// Check minimum length for security
	if len(password) < 8 {
		return "", fluxerrors.ErrPasswordTooShort
	}

	// Prehash long password due to bcrypt input limit.
	if len(password) > 72 {
		hasher := sha256.New()
		hasher.Write([]byte(password))
		password = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	}

	cost := 10

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func VerifyPassword(hashedPassword, password string) (matches bool, err error) {
	// For long passwords, apply the same pre-hash as during creation
	if len(password) > 72 {
		hasher := sha256.New()
		hasher.Write([]byte(password))
		password = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	}

	// CompareHashAndPassword returns nil on success, or an error on failure
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	matches = err == nil

	return
}
