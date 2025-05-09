package controller

import (
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
}

func NewS3CallbackController(bucketName string, videoSvc *service.VideoService) *S3CallbackController {
	return &S3CallbackController{
		bucketName: bucketName,
		vidSvc:     videoSvc,
	}
}

func (s *S3CallbackController) HandleVideoUploadEvent(c *gin.Context) {
	var event events.S3Event

	if err := c.BindJSON(&event); err != nil {
		response.Error(c, response.StatusBadRequest, "Malformed data", "Malformed data found.")
		return
	}

	// TODO: Add Concurrent handling of events
	for _, record := range event.Records {

		// Check if the event if object put and the bucket matches
		if strings.EqualFold(record.EventName, "s3:ObjectCreated:Put") && strings.EqualFold(record.S3.Bucket.Name, s.bucketName) {

			// Clean the object name to get the key with ID.
			videoSlug := strings.Replace(record.S3.Object.Key, fmt.Sprintf("%s/", s.bucketName), "", 1)

			err := s.vidSvc.UpdateUploadStatus(c.Request.Context(), videoSlug, model.UpdateVideoMeta{
				StoragePath: record.S3.Object.Key,
			})

			// TODO: Add service logger
			if err != nil {
				continue
			}

		}
	}

	response.Success(c, response.StatusOK, "", "")
}
