package repository

import (
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
)

type UserRepository struct {
	db *pgsql.PgSQL
}

func NewUserRepository(db *pgsql.PgSQL) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) CreateUser(user *model.User) (id model.UserID, err error) {

	userTable := tables.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	}

	result := u.db.DB.Create(&userTable)

	if result.Error != nil {
		err = result.Error
		return
	}

	id = model.UserID(userTable.ID.String())

	return
}
