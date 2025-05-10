package repository

import (
	"context"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/repository/pgsql"
	"fluxio-backend/pkg/repository/pgsql/tables"
	"fluxio-backend/pkg/utils"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VideoRepository struct {
	db *pgsql.PgSQL

	s3Client   *s3.S3
	bucketName string
}

type VideoRepositoryConfig struct {
	S3BucketName string
	S3Region     string
	S3AccessKey  string
	S3SecretKey  string
	S3Endpoint   string
}

func NewVideoRepository(db *pgsql.PgSQL, cfg VideoRepositoryConfig) *VideoRepository {

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

	return &VideoRepository{db: db, s3Client: s3Client, bucketName: cfg.S3BucketName}
}

func (r *VideoRepository) CreateVideoMeta(ctx context.Context, videoMeta model.Video) (video model.Video, err error) {

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
func (r *VideoRepository) IncrementVideoRetryCount(ctx context.Context, videoID model.VideoID) (err error) {
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

// UpdateMeta updates video metadata with the provided parameters
func (r *VideoRepository) UpdateMeta(ctx context.Context, id model.VideoID, status model.VideoStatus, params model.UpdateVideoMeta) (err error) {
	// Parse the VideoID to UUID
	uuid, err := uuid.Parse(id.String())
	if err != nil {
		return fluxerrors.ErrInvalidVideoID
	}

	// Validate the video status
	if !status.IsAcceptable() {
		return fluxerrors.ErrInvalidVideoStatus
	}

	// Create a map of fields to update based on the params struct
	updateData := r.buildUpdateDataMap(status, params)

	// Execute the update
	tx := r.db.DB.WithContext(ctx).Model(&tables.Video{}).Where("id = ?", uuid).Updates(updateData)

	if tx.Error != nil {
		return fluxerrors.ErrVideoMetaUpdateFailed
	}

	if tx.RowsAffected == 0 {
		return fluxerrors.ErrVideoNotFound
	}

	return nil
}

// buildUpdateDataMap is a private helper method that constructs the update data map
// from the provided status and UpdateVideoMeta parameters
func (r *VideoRepository) buildUpdateDataMap(status model.VideoStatus, params model.UpdateVideoMeta) map[string]interface{} {
	// Initialize with status and reset retry count
	updateData := map[string]interface{}{}

	// Add StoragePath from params
	if !strings.EqualFold(params.StoragePath, "") {
		updateData["storage_path"] = params.StoragePath
	}

	// Add other fields from params, only if they're not zero values
	// Title
	if !strings.EqualFold(params.Title, "") {
		updateData["title"] = params.Title
	}

	// Description
	if !strings.EqualFold(params.Description, "") {
		updateData["description"] = params.Description
	}

	// ParentID
	if params.ParentID != nil {
		updateData["parent_id"] = params.ParentID.String()
	}

	// Width
	if params.Width > 0 {
		updateData["width"] = params.Width
	}

	// Height
	if params.Height > 0 {
		updateData["height"] = params.Height
	}

	// Format
	if !strings.EqualFold(params.Format, "") {
		updateData["format"] = params.Format
	}

	// Status
	if !strings.EqualFold(params.Status.String(), "") {
		updateData["status"] = params.Status.String()
	}

	// Length
	if params.Length > 0 {
		updateData["length"] = params.Length
	}

	// AudioSampleRate
	if params.AudioSampleRate > 0 {
		updateData["audio_sample_rate"] = params.AudioSampleRate
	}

	// AudioCodec
	if !strings.EqualFold(params.AudioCodec, "") {
		updateData["audio_codec"] = params.AudioCodec
	}

	// IsFeatured - this is a boolean, so we always include it
	updateData["is_featured"] = params.IsFeatured

	// Visibility - only update if valid
	if params.Visibility.IsAcceptable() {
		updateData["visibility"] = params.Visibility.String()
	}

	// Slug
	if !strings.EqualFold(params.Slug, "") {
		updateData["slug"] = params.Slug
	}

	// Size
	if params.Size > 0 {
		updateData["size"] = params.Size
	}

	// Language
	if !strings.EqualFold(params.Language, "") {
		updateData["language"] = params.Language
	}

	return updateData
}

func (r *VideoRepository) GetVideoByID(ctx context.Context, id model.VideoID) (video model.Video, err error) {
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

func (r *VideoRepository) GetVideoBySlug(ctx context.Context, slug string) (video model.Video, err error) {

	data := &tables.Video{}

	tx := r.db.DB.First(data, "slug = ?", slug)

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

func (r *VideoRepository) CheckVideoExistsByID(ctx context.Context, id model.VideoID) (exists bool, err error) {
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
func (r *VideoRepository) GetProcessingDetailsBySlug(ctx context.Context, slug string) (video model.Video, err error) {

	if strings.EqualFold(slug, "") {
		err = fluxerrors.ErrInvalidVideoSlug
		return
	}

	tableVid := tables.Video{}

	tx := r.db.DB.Model(&tables.Video{}).Select("id, title, is_featured, storage_path, status, created_at, updated_at, retry_count").Where("slug = ?", slug).Find(&tableVid)

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
