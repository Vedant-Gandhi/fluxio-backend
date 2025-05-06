package routes

import (
	"fluxio-backend/pkg/transport/http/controller"
	"fluxio-backend/pkg/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type AWSCallbackRouter struct {
	s3CbHandler *controller.S3CallbackController
	middleware  *middleware.Middleware
}

func NewAWSCallbackRouter(s3CbController *controller.S3CallbackController, middleware *middleware.Middleware) *AWSCallbackRouter {
	return &AWSCallbackRouter{
		s3CbHandler: s3CbController,
		middleware:  middleware,
	}
}

func (r *AWSCallbackRouter) RegisterRoutes(router *gin.Engine) {
	AwsGroup := router.Group("/api/v1/aws-cb")
	{
		AwsGroup.POST("/s3/upload-object", r.s3CbHandler.HandleVideoUploadEvent)

	}
}
