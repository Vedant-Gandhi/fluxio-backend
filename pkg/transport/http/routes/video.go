package routes

import (
	"fluxio-backend/pkg/transport/http/controller"
	"fluxio-backend/pkg/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type VideoRouter struct {
	VideoController *controller.VideoController
	middleware      *middleware.Middleware
}

func NewVideoRouter(VideoController *controller.VideoController, middleware *middleware.Middleware) *VideoRouter {
	return &VideoRouter{
		VideoController: VideoController,
		middleware:      middleware,
	}
}

// RegisterRoutes registers all Video-related routes
func (r *VideoRouter) RegisterRoutes(router *gin.Engine) {
	VideoGroup := router.Group("/api/v1/video")
	{
		VideoGroup.POST("/upload-init", r.middleware.Auth.Add(), r.VideoController.GenerateVideUploadEntry)

	}
}
