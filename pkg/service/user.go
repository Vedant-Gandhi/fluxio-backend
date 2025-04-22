package service

import (
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/fluxcrypto"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository"
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
