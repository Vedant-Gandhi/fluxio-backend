package repository

import (
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
	"fluxio-backend/pkg/utils"
	"strings"

	"gorm.io/gorm"
)

type VideoMetaRepository struct {
	db *pgsql.PgSQL
}

func NewVideoMetaRepository(db *pgsql.PgSQL) *VideoMetaRepository {

	return &VideoMetaRepository{db: db}
}

func (r *VideoMetaRepository) CreateVideoMeta(videoMeta model.Video) (video model.Video, err error) {

	if strings.EqualFold(videoMeta.Visibility.String(), "") {
		videoMeta.Visibility = model.VideoVisibilityPublic
	}

	slug := utils.CreateURLSafeVideoSlug(videoMeta.Title)
	if strings.EqualFold(slug, "") {
		err = fluxerrors.ErrFailedToGenerateVideoSlug
		return
	}

	vidTable := &tables.Video{
		Title:      videoMeta.Title,
		UserID:     videoMeta.UserID,
		Status:     model.VideoStatusPending.String(),
		Visibility: videoMeta.Visibility.String(),
		Slug:       slug,
	}

	tx := r.db.DB.Create(vidTable)

	err = tx.Error

	if err != nil {
		if err == gorm.ErrDuplicatedKey {
			if strings.Contains(err.Error(), "video_title_key") {
				err = fluxerrors.ErrDuplicateVideoTitle
				return
			}

			err = fluxerrors.ErrVideoAlreadyExists

			return
		}

		err = fluxerrors.ErrVideoCreationFailed
		return

	}

	return
}
