package controller

import (
	"context"
	"fluxio-backend/pkg/common/schema"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http/response"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
)

type S3CallbackController struct {
	bucketName string
	vidSvc     *service.VideoService

	l schema.Logger
}

func NewS3CallbackController(bucketName string, videoSvc *service.VideoService, logger schema.Logger) *S3CallbackController {
	return &S3CallbackController{
		bucketName: bucketName,
		vidSvc:     videoSvc,

		l: logger,
	}
}

func (s *S3CallbackController) HandleVideoUploadEvent(c *gin.Context) {
	logger := s.l.WithField(map[string]interface{}{
		"handler": "S3Callback",
		"bucket":  s.bucketName,
	})

	logger.Info("Received S3 event notification")

	var event events.S3Event

	if err := c.BindJSON(&event); err != nil {
		logger.Error("Failed to parse S3 event", err)
		response.Error(c, response.StatusBadRequest, "Malformed data", "Malformed data found.")
		return
	}

	logger.Info("Processing S3 event", "record_count", len(event.Records))

	// TODO: Add Concurrent handling of events
	for _, record := range event.Records {
		recordLogger := logger.With("event_name", record.EventName)

		// Check if the event if object put and the bucket matches
		if strings.EqualFold(record.EventName, "s3:ObjectCreated:Put") && strings.EqualFold(record.S3.Bucket.Name, s.bucketName) {
			// Clean the object name to get the key with ID.
			videoSlug := strings.Replace(record.S3.Object.Key, fmt.Sprintf("%s/", s.bucketName), "", 1)
			objectSize := record.S3.Object.Size

			recordLogger = recordLogger.With("video_slug", videoSlug).With("object_size", objectSize)
			recordLogger.Info("Processing video upload")

			err := s.vidSvc.UpdateUploadStatus(c.Request.Context(), videoSlug, model.UpdateVideoMeta{
				StoragePath: record.S3.Object.Key,
			})

			if err != nil {
				recordLogger.Error("Failed to update upload status", err)
				continue
			}

			recordLogger.Info("Starting asynchronous video processing")

			// Create a background context for the goroutine
			bgCtx := context.Background()

			// Use a goroutine for async processing
			go func(ctx context.Context, slug string) {
				processingLogger := s.l.With("video_slug", slug).With("processing", "async")
				processingLogger.Info("Beginning post-upload processing")

				err := s.vidSvc.PerformPostUploadProcessing(ctx, slug)
				if err != nil {
					processingLogger.Error("Post-upload processing failed", err)
					return
				}

				processingLogger.Info("Post-upload processing completed successfully")
			}(bgCtx, videoSlug)
		} else {
			recordLogger.Info("Skipping non-matching event",
				"event_bucket", record.S3.Bucket.Name,
				"expected_bucket", s.bucketName)
		}
	}

	logger.Info("S3 event handling completed")
	response.Success(c, response.StatusOK, "", "")
}
