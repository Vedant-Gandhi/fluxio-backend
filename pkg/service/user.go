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
	logger   schema.Logger
}

func NewUserService(repo *repository.UserRepository, jwtsvc *JWTService, logger schema.Logger) *UserService {
	return &UserService{
		repo:     repo,
		jService: jwtsvc,
		logger:   logger,
	}
}

func (s *UserService) CreateUser(user model.User, rawPassword string) (id model.UserID, token string, err error) {
	logger := s.logger.With("email", user.Email)

	hashedPass, err := fluxcrypto.HashPassword(rawPassword)

	if err != nil {
		logger.Error("Failed to create a user password hash", err)
		err = fluxerrors.ErrPasswordFailed
		return
	}

	user.Password = hashedPass

	id, err = s.repo.CreateUser(user)

	if err != nil {
		logger.Error("Failed to create a new user", err)
		// Only allow the username and email errors to populate.
		if err == fluxerrors.ErrUsernameExists || err == fluxerrors.ErrEmailExists {
			return
		}
		fmt.Printf("Error creating user: %v", err)
		err = fluxerrors.ErrUserCreationFailed
		return
	}

	logger.Info("User created successfully", "user_id", id.String())

	// Suppress the error if token creation fails.
	token, tokenErr := s.jService.GenerateToken(model.JWTTokenClaims{
		UserID: user.ID.String(),
		Sub:    "user",
	}, time.Now().Add(time.Duration(TOKEN_EXPIRY)*time.Second))

	if tokenErr != nil {
		logger.Error("Token generation failed", tokenErr)
	}

	return
}

func (s *UserService) Login(userData model.User) (user model.User, token string, err error) {
	logger := s.logger
	if userData.Email != "" {
		logger = logger.With("email", userData.Email)
	} else {
		logger = logger.With("username", userData.Username)
	}

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
			if err != fluxerrors.ErrUserNotFound {
				logger.Error("Failed to get a user by username", err)
			}
			return
		}
	}

	// Verify by email if username not found
	if !emailEmpty && strings.EqualFold(user.ID.String(), "") {
		user, err = s.repo.GetUserByEmail(userData.Email)

		if err != nil {
			if err != fluxerrors.ErrUserNotFound {
				logger.Error("Failed to get a user by email", err)
			}
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
		logger.Error("Failed to generate token", err)
		err = fluxerrors.ErrFailedToCreateToken

		// Reset the user object to prevent leaking user data.
		user = model.User{}
		return
	}

	logger.Info("User authenticated", "user_id", user.ID.String())
	return
}

func (s *UserService) GetUserByID(userID model.UserID) (user model.User, err error) {
	logger := s.logger.With("user_id", userID.String())

	user, err = s.repo.GetUserByID(userID)

	if err != nil {
		if err == fluxerrors.ErrUserNotFound || err == fluxerrors.ErrInvalidUserID {
			return
		}
		logger.Error("Failed to fetch user", err)
		err = fluxerrors.ErrUnknown
		return
	}

	return
}
