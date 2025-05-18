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

	var video model.Video

	if err := c.ShouldBindJSON(&video); err != nil {
		response.Error(c, response.StatusBadRequest, "Invalid request payload", "The payload is not valid.")
		return
	}
	video, uploadURL, err := v.videoService.AddVideo(c, video)
	if err != nil {
		if err == fluxerrors.ErrDuplicateVideoTitle {
			response.Error(c, response.StatusConflict, response.MsgDuplicateVideoTitle, err.Error())
			return
		}
		response.Error(c, response.StatusUnprocessableEntity, response.MsgVideoCreationFailed, err.Error())
		return
	}
	response.Success(c, response.StatusCreated, "Video entry created successfully", gin.H{
		"video":      video,
		"upload_url": uploadURL.String(),
	})

}
