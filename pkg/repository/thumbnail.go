package repository

import (
	"context"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql/tables"

	"github.com/google/uuid"
)

func (v *VideoRepository) CreateThumbnail(ctx context.Context, thumbnail model.Thumbnail) (id model.ThumbnailID, err error) {

	if thumbnail.Width == 0 || thumbnail.Height == 0 {
		err = fluxerrors.ErrInvalidThumbnailDimensions
		return
	}

	parsedVidId, err := uuid.Parse(thumbnail.VideoID.String())
	if err != nil {
		err = fluxerrors.ErrInvalidVideoID
		return
	}

	insertData := tables.Thumbnail{
		VideoID:     parsedVidId,
		Width:       thumbnail.Width,
		Height:      thumbnail.Height,
		Size:        thumbnail.Size,
		Format:      thumbnail.Format,
		StoragePath: thumbnail.StoragePath,
		TimeStamp:   thumbnail.TimeStamp,
	}

	tx := v.db.DB.Create(&insertData)

	if tx.Error != nil {
		err = fluxerrors.ErrThumbnailCreationFailed
		return
	}

	id = model.ThumbnailID(insertData.ID.String())

	return
}
