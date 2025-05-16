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

func (v *VideoRepository) GenerateUnProcessedVideoUploadURL(ctx context.Context, id model.VideoID, slug string) (url *url.URL, err error) {

	path := v.generateVideoFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimLeft(path, fmt.Sprintf("%s/", v.rawVidBketName))

	s3Request, _ := v.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(v.rawVidBketName),
		Key:         aws.String(path),
		ContentType: aws.String("application/octet-stream"),
	})

	rawURL, err := s3Request.Presign(constants.PreSignedVidUploadURLExpireTime)

	if err != nil {
		err = fluxerrors.ErrVideoURLGenerationFailed
		return
	}

	url, _ = url.Parse(rawURL)

	return
}

func (v *VideoRepository) GetUnProcessedVideoDownloadURL(ctx context.Context, slug string) (url *url.URL, err error) {

	path := v.generateVideoFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimLeft(path, fmt.Sprintf("%s/", v.rawVidBketName))

	s3Request, _ := v.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(v.rawVidBketName),
		Key:    aws.String(path),
	})

	rawURL, err := s3Request.Presign(constants.PreSignedVidTempDownloadURLExpireTime)

	if err != nil {
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
