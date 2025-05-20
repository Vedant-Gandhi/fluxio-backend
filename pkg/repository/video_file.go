package repository

import (
	"context"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/utils"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (v *VideoRepository) GenerateUnProcessedVideoUploadURL(ctx context.Context, id model.VideoID, slug string, mimeType string) (url *url.URL, err error) {
	logger := v.l.With("video_id", id.String())
	path := v.generateVideoFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimPrefix(path, fmt.Sprintf("%s/", v.rawVidBketName))

	s3Request, _ := v.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(v.rawVidBketName),
		Key:         aws.String(path),
		ContentType: aws.String("video/mp4"),
	})

	rawURL, err := s3Request.Presign(constants.PreSignedVidUploadURLExpireTime)

	if err != nil {
		logger.Error("Failed to create a presigned URL for video upload", err)
		err = fluxerrors.ErrVideoURLGenerationFailed
		return
	}

	url, _ = url.Parse(rawURL)

	return
}

func (v *VideoRepository) GetUnProcessedVideoDownloadURL(ctx context.Context, slug string) (url *url.URL, err error) {
	logger := v.l.With("video_slug", slug)
	path := v.generateVideoFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimPrefix(path, fmt.Sprintf("%s/", v.rawVidBketName))

	s3Request, _ := v.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(v.rawVidBketName),
		Key:    aws.String(path),
	})

	rawURL, err := s3Request.Presign(constants.PreSignedVidTempDownloadURLExpireTime)

	if err != nil {
		logger.Error("Failed to create a presigned URL for video download", err)
		err = fluxerrors.ErrVideoURLGenerationFailed
		return
	}

	url, _ = url.Parse(rawURL)

	return
}

func (v *VideoRepository) generateVideoFileS3Path(slug string) string {
	return strings.TrimRight(slug, "/")
}

func (v *VideoRepository) generateThumbnailFileS3Path(id model.VideoID, timestamp uint64, extension string) string {
	path := utils.CreateURLSafeThumbnailFileName(id.String(), fmt.Sprint(timestamp))

	// Add file extension to the path.
	path = fmt.Sprintf("%s.%s", path, extension)
	return strings.TrimSpace(path)
}
