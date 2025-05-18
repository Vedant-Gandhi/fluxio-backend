package repository

import (
	"fluxio-backend/pkg/common/schema"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db     *pgsql.PgSQL
	logger schema.Logger
}

func NewUserRepository(db *pgsql.PgSQL, logger schema.Logger) *UserRepository {
	return &UserRepository{db: db, logger: logger}
}

func (u *UserRepository) CreateUser(user model.User) (id model.UserID, err error) {
	logger := u.logger.With("email", user.Email)
	userTable := tables.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	result := u.db.DB.Create(&userTable)

	if result.Error != nil {
		logger.Error("Error when creating a user", result.Error)

		// Check for duplicate key violation
		if result.Error == gorm.ErrDuplicatedKey || strings.Contains(result.Error.Error(), "uni_users_username") || strings.Contains(result.Error.Error(), "uni_users_email") {

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

func (u *UserRepository) CheckUserExists(id model.UserID) (exists bool, err error) {
	logger := u.logger.With("user_id", id.String())
	uuid, err := uuid.Parse(id.String())

	if err != nil {
		err = fluxerrors.ErrInvalidUserID
		return
	}

	user := &tables.User{}

	result := u.db.DB.First(user, " id = ?", uuid)

	if result.Error != nil {
		logger.Error("Error when checking user exists", result.Error)
		err = result.Error
	}

	// Ensure the user has not been deleted.
	exists = result.Error != gorm.ErrRecordNotFound

	return
}

func (u *UserRepository) GetUserByID(id model.UserID) (user model.User, err error) {
	logger := u.logger.With("user_id", id.String())
	uuid, err := uuid.Parse(id.String())

	if err != nil {
		err = fluxerrors.ErrInvalidUserID
		return
	}

	userTable := &tables.User{}

	result := u.db.DB.First(userTable, " id = ?", uuid)

	if result.Error != nil {
		logger.Error("Error when getting a user by ID", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			err = fluxerrors.ErrUserNotFound
			return
		}
		err = result.Error
		return
	}

	return
}

func (u *UserRepository) GetUserByUsername(username string) (user model.User, err error) {
	logger := u.logger.With("username", username)

	userTable := &tables.User{}

	result := u.db.DB.First(userTable, " username = ?", username)

	if result.Error != nil {
		logger.Error("Error when getting a user by username", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			err = fluxerrors.ErrUserNotFound
			return
		}
		err = result.Error
		return
	}

	user = model.User{
		ID:            model.UserID(userTable.ID.String()),
		Username:      userTable.Username,
		Email:         userTable.Email,
		Password:      userTable.Password,
		UpdatedAt:     userTable.UpdatedAt,
		CreatedAt:     userTable.CreatedAt,
		IsBlackListed: userTable.IsBlackListed,
	}

	return
}

func (u *UserRepository) GetUserByEmail(email string) (user model.User, err error) {
	logger := u.logger.With("email", email)
	userTable := &tables.User{}

	result := u.db.DB.First(userTable, " email = ?", email)

	if result.Error != nil {
		logger.Error("Error when getting a user by username", result.Error)
		if result.Error == gorm.ErrRecordNotFound {
			err = fluxerrors.ErrUserNotFound
			return
		}
		err = result.Error
		return
	}

	user = model.User{
		ID:            model.UserID(userTable.ID.String()),
		Username:      userTable.Username,
		Email:         userTable.Email,
		Password:      userTable.Password,
		UpdatedAt:     userTable.UpdatedAt,
		CreatedAt:     userTable.CreatedAt,
		IsBlackListed: userTable.IsBlackListed,
	}

	return
}
