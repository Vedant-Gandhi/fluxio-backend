package routes

import (
	"fluxio-backend/pkg/transport/http/controller"
	"fluxio-backend/pkg/transport/http/middleware"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authController *controller.AuthController
	middleware     *middleware.Middleware
}

func NewAuthRouter(authController *controller.AuthController, middleware *middleware.Middleware) *AuthRouter {
	return &AuthRouter{
		authController: authController,
		middleware:     middleware,
	}
}

// RegisterRoutes registers all auth-related routes
func (r *AuthRouter) RegisterRoutes(router *gin.Engine) {
	authGroup := router.Group("/api/v1/auth")
	{
		authGroup.POST("/register", r.authController.RegisterUser)
		authGroup.POST("/login", r.authController.LoginUser)

	}
}
