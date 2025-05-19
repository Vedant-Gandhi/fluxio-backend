package controller

import (
	"fluxio-backend/pkg/common/schema"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http/response"

	"github.com/gin-gonic/gin"
)

type VideoController struct {
	videoService *service.VideoService

	l schema.Logger
}

func NewVideoController(videoService *service.VideoService, logger schema.Logger) *VideoController {
	return &VideoController{
		videoService: videoService,
		l:            logger,
	}
}

func (v *VideoController) CreateNewVideo(c *gin.Context) {
	logger := v.l.WithField(map[string]interface{}{
		"client_ip": c.ClientIP(),
		"path":      c.FullPath(),
		"method":    c.Request.Method,
	})

	logger.Info("Handling video creation request")

	var video model.Video

	if err := c.ShouldBindJSON(&video); err != nil {
		logger.Warn("Invalid video creation payload found")
		logger.Debug("Invalid video creation payload", err)
		response.Error(c, response.StatusBadRequest, "Invalid request payload", "The payload is not valid.")
		return
	}

	logger = logger.With("title", video.Title)

	video, uploadURL, err := v.videoService.AddVideo(c, video)
	if err != nil {
		if err == fluxerrors.ErrDuplicateVideoTitle {
			logger.Info("Video creation failed - duplicate title")
			response.Error(c, response.StatusConflict, response.MsgDuplicateVideoTitle, err.Error())
			return
		}

		logger.Error("Video creation failed", err)
		response.Error(c, response.StatusUnprocessableEntity, response.MsgVideoCreationFailed, err.Error())
		return
	}

	logger = logger.With("video_id", video.ID.String()).With("slug", video.Slug)
	logger.Info("Video entry created successfully")

	response.Success(c, response.StatusCreated, "Video entry created successfully", gin.H{
		"video":      video,
		"upload_url": uploadURL.String(),
	})
}
