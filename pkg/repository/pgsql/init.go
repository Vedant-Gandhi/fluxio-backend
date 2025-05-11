package pgsql

import (
	"fluxio-backend/pkg/repository/pgsql/tables"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type PgSQLConfig struct {
	URL string
}

type PgSQL struct {
	DB *gorm.DB
}

func NewPgSQL(cfg PgSQLConfig) (*PgSQL, error) {

	// Configure GORM
	gormConfig := &gorm.Config{}

	// Enable SQL logging in development
	if os.Getenv("GO_ENV") == "development" {
		gormConfig.Logger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.URL), gormConfig)

	if err != nil {
		return nil, err
	}

	// Automatically migrate the schema, keeping your database up to date.
	db.AutoMigrate(&tables.User{})
	db.AutoMigrate(&tables.Video{})
	db.AutoMigrate(&tables.Thumbnail{})

	return &PgSQL{
		DB: db,
	}, nil

}
