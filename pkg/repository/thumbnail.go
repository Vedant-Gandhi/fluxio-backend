package repository

import (
	"context"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql/tables"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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

func (v *VideoRepository) GenerateThumbnailUploadURL(ctx context.Context, id model.VideoID, timestamp uint64, extension string) (url *url.URL, err error) {

	path := v.generateThumbnailFileS3Path(id, timestamp, extension)

	if strings.EqualFold(path, "") {
		err = fluxerrors.ErrThumbnailURLGenerationFailed
		return
	}
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimLeft(path, fmt.Sprintf("%s/", v.thumbnailBucketName))

	s3Request, _ := v.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(v.thumbnailBucketName),
		Key:         aws.String(path),
		ContentType: aws.String(fmt.Sprintf("image/%s", extension)),
	})

	rawURL, err := s3Request.Presign(constants.PreSignedVidUploadURLExpireTime)

	if err != nil {
		err = fluxerrors.ErrThumbnailURLGenerationFailed
		return
	}

	url, _ = url.Parse(rawURL)

	return
}
