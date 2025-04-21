package pgsql

import (
	"fluxio-backend/pkg/repository/pgsql/tables"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PgSQLConfig struct {
	URL string
}

type PgSQL struct {
	db *gorm.DB
}

func NewPgSQL(cfg PgSQLConfig) (*PgSQL, error) {

	db, err := gorm.Open(postgres.Open(cfg.url), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema, keeping your database up to date.
	db.AutoMigrate(&tables.User{})

	return &PgSQL{
		db: db,
	}, nil

}
