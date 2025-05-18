package service

import (
	"fluxio-backend/pkg/common/schema"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/fluxcrypto"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"fmt"
	"strings"
	"time"
)

const TOKEN_EXPIRY = uint(8 * 60) // 8 hours

type UserService struct {
	repo     *repository.UserRepository
	jService *JWTService
	l        schema.Logger
}

func NewUserService(repo *repository.UserRepository, jwtsvc *JWTService, logger schema.Logger) *UserService {
	return &UserService{
		repo:     repo,
		jService: jwtsvc,
		l:        logger,
	}
}

func (s *UserService) CreateUser(user model.User, rawPassword string) (id model.UserID, token string, err error) {

	hashedPass, err := fluxcrypto.HashPassword(rawPassword)

	if err != nil {
		err = fluxerrors.ErrPasswordFailed
		return
	}

	user.Password = hashedPass

	id, err = s.repo.CreateUser(user)

	if err != nil {

		// Only allow the username and email errors to populate.
		if err == fluxerrors.ErrUsernameExists || err == fluxerrors.ErrEmailExists {
			return
		}
		fmt.Printf("Error creating user: %v", err)
		err = fluxerrors.ErrUserCreationFailed
		return
	}

	// Suppress the error if token creation fails.
	token, _ = s.jService.GenerateToken(model.JWTTokenClaims{
		UserID: user.ID.String(),
		Sub:    "user",
	}, time.Now().Add(time.Duration(TOKEN_EXPIRY)*time.Second))

	return
}

func (s *UserService) Login(userData model.User) (user model.User, token string, err error) {

	usernameEmpty := strings.EqualFold(userData.Username, "")
	emailEmpty := strings.EqualFold(userData.Email, "")

	userInputPassword := userData.Password

	if usernameEmpty && emailEmpty {
		err = fluxerrors.ErrInvalidCredentials
		return
	}

	// Verify by username first
	if !usernameEmpty {
		user, err = s.repo.GetUserByUsername(userData.Username)

		if err != nil {
			return
		}

	}

	// Verify by email if username not found
	if !emailEmpty && strings.EqualFold(user.ID.String(), "") {
		user, err = s.repo.GetUserByEmail(userData.Email)

		if err != nil {
			return
		}

	}

	if strings.EqualFold(user.ID.String(), "") {
		err = fluxerrors.ErrUserNotFound
		return
	}

	// Check password if user found.
	matches, _ := fluxcrypto.VerifyPassword(user.Password, userInputPassword)
	if !matches {
		err = fluxerrors.ErrInvalidCredentials
		return
	}

	token, err = s.jService.GenerateToken(model.JWTTokenClaims{
		UserID: user.ID.String(),
		Sub:    "user",
	}, time.Now().Add(time.Duration(TOKEN_EXPIRY)*time.Second))

	if err != nil {
		err = fluxerrors.ErrFailedToCreateToken

		// Reset the user object to prevent leaking user data.
		user = model.User{}
	}

	return

}

func (s *UserService) GetUserByID(userID model.UserID) (user model.User, err error) {

	user, err = s.repo.GetUserByID(userID)

	if err != nil {
		if err == fluxerrors.ErrUserNotFound || err == fluxerrors.ErrInvalidUserID {
			return
		}
		err = fluxerrors.ErrUnknown
		return

	}

	return
}
