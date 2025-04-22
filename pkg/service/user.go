package service

import (
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/fluxcrypto"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
	"strings"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(user model.User, rawPassword string) (id model.UserID, err error) {

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
		err = fluxerrors.ErrUserCreationFailed
		return
	}

	return
}

func (s *UserService) Login(userData model.User, rawPassword string) (user model.User, err error) {

	usernameEmpty := strings.EqualFold(userData.Username, "")
	emailEmpty := strings.EqualFold(userData.Email, "")

	if usernameEmpty && emailEmpty {
		err = fluxerrors.ErrInvalidCredentials
		return
	}

	// Verify by username first
	if !usernameEmpty {
		user, err = s.repo.GetUserByUsername(userData.Username)

		if err == nil {
			matches, err := fluxcrypto.VerifyPassword(user.Password, rawPassword)
			if !matches {
				err = fluxerrors.ErrInvalidCredentials
			}

			return user, err
		}
	}

	// Verify by email if username not found
	if !emailEmpty {
		user, err = s.repo.GetUserByEmail(userData.Email)

		if err == nil {
			matches, err := fluxcrypto.VerifyPassword(user.Password, rawPassword)
			if !matches {
				err = fluxerrors.ErrInvalidCredentials
			}

			return user, err
		}

	}

	return
}
