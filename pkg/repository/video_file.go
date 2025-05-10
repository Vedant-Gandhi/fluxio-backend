package repository

import (
	"context"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func (v *VideoRepository) GenerateVideoUploadURL(ctx context.Context, id model.VideoID, slug string) (url *url.URL, err error) {

	path := v.generateFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimLeft(path, fmt.Sprintf("%s/", v.bucketName))

	s3Request, _ := v.s3Client.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(v.bucketName),
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

func (v *VideoRepository) GetVideoTemporaryDownloadURL(ctx context.Context, slug string) (url *url.URL, err error) {

	path := v.generateFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimLeft(path, fmt.Sprintf("%s/", v.bucketName))

	s3Request, _ := v.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(v.bucketName),
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

func (v *VideoRepository) generateFileS3Path(slug string) string {
	return fmt.Sprintf("%s/%s", v.bucketName, strings.TrimRight(slug, "/"))
}
