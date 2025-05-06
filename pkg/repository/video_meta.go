package repository

import (
	"context"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
	"fluxio-backend/pkg/utils"
	"strings"

	"github.com/google/uuid"
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

func (r *VideoMetaRepository) UpdateProcessingDetails(ctx context.Context, id model.VideoID, status model.VideoStatus, storagePath string) (err error) {

	if strings.EqualFold(storagePath, "") {
		err = fluxerrors.ErrMalformedStoragePath
		return
	}
	uuid, err := uuid.Parse(id.String())

	if err != nil {
		err = fluxerrors.ErrInvalidVideoID
		return
	}

	if !status.IsAcceptable() {
		err = fluxerrors.ErrInvalidVideoStatus
		return
	}

	updateData := map[string]any{
		"storage_path": storagePath,
		"status":       status.String(),
		"retry_count":  0,
	}

	tx := r.db.DB.WithContext(ctx).Model(&tables.Video{}).Where("id = ?", uuid).Updates(updateData)

	if tx.Error != nil {
		err = fluxerrors.ErrVideoMetaUpdateFailed
		return
	}

	if tx.RowsAffected == 0 {
		err = fluxerrors.ErrVideoNotFound
		return
	}

	return

}

func (r *VideoMetaRepository) GetVideoByID(ctx context.Context, id model.VideoID) (video model.Video, err error) {
	uuid, err := uuid.Parse(id.String())

	if err != nil {
		err = fluxerrors.ErrInvalidVideoID
		return
	}

	data := &tables.Video{}

	tx := r.db.DB.First(data, "id = ?", uuid)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			err = fluxerrors.ErrVideoNotFound
			return
		}

		return
	}

	video = model.Video{
		ID:         model.VideoID(data.ID.String()),
		Title:      data.Title,
		UserID:     data.UserID,
		Status:     model.VideoStatus(data.Status),
		Visibility: model.VideoVisibility(data.Visibility),
		Slug:       data.Slug,
		RetryCount: data.RetryCount,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
		IsFeatured: data.IsFeatured,
	}

	return
}

func (r *VideoMetaRepository) CheckVideoExistsByID(ctx context.Context, id model.VideoID) (exists bool, err error) {
	uuid, err := uuid.Parse(id.String())

	if err != nil {
		err = fluxerrors.ErrInvalidVideoID
		return
	}

	tx := r.db.DB.Model(&tables.Video{}).Select("count(*) > 0").Where("id = ?", uuid).Find(&exists)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			err = nil
			return
		}

		return
	}

	return
}

// Returns the details required for video processing.
func (r *VideoMetaRepository) GetProcessingDetailsByID(ctx context.Context, id model.VideoID) (video model.Video, err error) {
	uuid, err := uuid.Parse(id.String())

	if err != nil {
		err = fluxerrors.ErrInvalidVideoID
		return
	}

	tableVid := tables.Video{}

	tx := r.db.DB.Model(&tables.Video{}).Select("id, title, is_featured, storage_path, status, created_at, updated_at, retry_count").Where("id = ?", uuid).Find(&tableVid)

	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			err = fluxerrors.ErrVideoNotFound
			return
		}

		return
	}

	video = model.Video{
		ID:          model.VideoID(tableVid.ID.String()),
		IsFeatured:  tableVid.IsFeatured,
		Status:      model.VideoStatus(tableVid.Status),
		StoragePath: tableVid.StoragePath,
		CreatedAt:   tableVid.CreatedAt,
		UpdatedAt:   tableVid.UpdatedAt,
		RetryCount:  tableVid.RetryCount,
	}

	return
}
