package routes

import (
	"fluxio-backend/pkg/transport/http/controller"

	"github.com/gin-gonic/gin"
)

type VideoRouter struct {
	VideoController *controller.VideoController
}

func NewVideoRouter(VideoController *controller.VideoController) *VideoRouter {
	return &VideoRouter{
		VideoController: VideoController,
	}
}

// RegisterRoutes registers all Video-related routes
func (r *VideoRouter) RegisterRoutes(router *gin.Engine) {
	VideoGroup := router.Group("/api/v1/video")
	{
		VideoGroup.POST("/upload", r.VideoController.GenerateVideUploadEntry)

	}
}
