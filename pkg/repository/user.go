package repository

import "fluxio-backend/pkg/repository/pgsql"

type UserRepository struct {
	db *pgsql.PgSQL
}

func NewUserRepository(db *pgsql.PgSQL) *UserRepository {
	return &UserRepository{db: db}
}
