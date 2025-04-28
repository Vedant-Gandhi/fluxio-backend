package repository

import (
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
)

type VideoMetaRepository struct {
	db *pgsql.PgSQL
}

func NewVideoMetaRepository(db *pgsql.PgSQL) *VideoMetaRepository {
	return &VideoMetaRepository{db: db}
}

func (r *VideoMetaRepository) CreateVideoMeta(videoMeta model.Video) (video model.Video, err error) {
	vidTable := &tables.Video{
		Title:  videoMeta.Title,
		UserID: videoMeta.UserID,
	}
	return
}
