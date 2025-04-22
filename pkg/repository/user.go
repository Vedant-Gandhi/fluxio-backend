package repository

import (
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
	"strings"
)

type UserRepository struct {
	db *pgsql.PgSQL
}

func NewUserRepository(db *pgsql.PgSQL) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(user model.User) (id model.UserID, err error) {
	userTable := tables.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	result := u.db.DB.Create(&userTable)

	if result.Error != nil {
		// Check for duplicate key violation
		if strings.Contains(result.Error.Error(), "unique constraint") ||
			strings.Contains(result.Error.Error(), "duplicate key") {

			// Check which unique constraint was violated
			if strings.Contains(result.Error.Error(), "username") {
				return "", fluxerrors.ErrUsernameExists
			} else if strings.Contains(result.Error.Error(), "email") {
				return "", fluxerrors.ErrEmailExists
			}

			return "", fluxerrors.ErrUserAlreadyExists
		}

		return "", result.Error
	}

	id = model.UserID(userTable.ID.String())
	return
}
