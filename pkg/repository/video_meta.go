package repository

import (
	"context"
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

func (r *VideoMetaRepository) CreateVideoMeta(ctx context.Context, videoMeta model.Video) (video model.Video, err error) {

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

	tx := r.db.DB.WithContext(ctx).Create(vidTable)

	err = tx.Error

	if err != nil {
		if err == gorm.ErrDuplicatedKey {

			err = fluxerrors.ErrVideoAlreadyExists

			return
		}

		if strings.Contains(err.Error(), "uni_videos_title") {
			err = fluxerrors.ErrDuplicateVideoTitle
			return
		}

		err = fluxerrors.ErrVideoCreationFailed
		return

	}

	video = model.Video{
		ID:         model.VideoID(vidTable.ID.String()),
		Title:      vidTable.Title,
		UserID:     vidTable.UserID,
		Status:     model.VideoStatus(vidTable.Status),
		Visibility: model.VideoVisibility(vidTable.Visibility),
		Slug:       vidTable.Slug,
		RetryCount: vidTable.RetryCount,
		CreatedAt:  vidTable.CreatedAt,
		UpdatedAt:  vidTable.UpdatedAt,
		IsFeatured: vidTable.IsFeatured,
	}

	if vidTable.DeletedAt.Valid {
		video.DeletedAt = &vidTable.DeletedAt.Time
	}

	return
}
func (r *VideoMetaRepository) IncrementVideoRetryCount(ctx context.Context, videoID model.VideoID) (err error) {
	vidTable := &tables.Video{}

	tx := r.db.DB.WithContext(ctx).Model(vidTable).
		Where("id = ?", videoID).
		Update("retry_count", gorm.Expr("retry_count + 1"))

	err = tx.Error

	if err != nil {
		return
	}

	if tx.RowsAffected == 0 {
		err = fluxerrors.ErrVideoNotFound
		return
	}

	return
}
