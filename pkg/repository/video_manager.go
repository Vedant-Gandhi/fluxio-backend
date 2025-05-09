package repository

import (
	"context"
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type VideoManagerRepositoryConfig struct {
	S3BucketName string
	S3Region     string
	S3AccessKey  string
	S3SecretKey  string
	S3Endpoint   string
}

type VideoManagerRepository struct {
	db         *pgsql.PgSQL
	awsS3      *s3.S3
	bucketName string
}

func NewVideoManagerRepository(db *pgsql.PgSQL, cfg VideoManagerRepositoryConfig) *VideoManagerRepository {

	url, _ := url.Parse(cfg.S3Endpoint)

	awsConfig := &aws.Config{
		Region:      aws.String(cfg.S3Region),
		Credentials: credentials.NewStaticCredentials(cfg.S3AccessKey, cfg.S3SecretKey, ""),
	}

	if !strings.EqualFold(url.Host, "") {

		awsConfig.Endpoint = aws.String(url.String())
		awsConfig.DisableSSL = aws.Bool(strings.EqualFold(url.Scheme, "http"))
		awsConfig.S3ForcePathStyle = aws.Bool(true)
	}

	awsSession, _ := session.NewSession(awsConfig)

	s3Client := s3.New(awsSession)

	return &VideoManagerRepository{db: db, awsS3: s3Client, bucketName: cfg.S3BucketName}
}

func (v *VideoManagerRepository) GenerateVideoUploadURL(ctx context.Context, id model.VideoID, slug string) (url *url.URL, err error) {

	path := v.generateFileS3Path(slug)
	// Remove the bucket name from the path to avoid double prefixing.
	path = strings.TrimLeft(path, fmt.Sprintf("%s/", v.bucketName))

	s3Request, _ := v.awsS3.PutObjectRequest(&s3.PutObjectInput{
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

func (v *VideoManagerRepository) GetVideoFileMeta(ctx context.Context, slug string) (err error) {

	return

}

func (v *VideoManagerRepository) generateFileS3Path(slug string) string {
	return fmt.Sprintf("%s/%s", v.bucketName, strings.TrimRight(slug, "/"))
}
