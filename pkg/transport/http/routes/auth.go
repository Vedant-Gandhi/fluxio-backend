package routes

import (
	"fluxio-backend/pkg/transport/http/controller"

	"github.com/gin-gonic/gin"
)

type AuthRouter struct {
	authController *controller.AuthController
}

func NewAuthRouter(authController *controller.AuthController) *AuthRouter {
	return &AuthRouter{
		authController: authController,
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
