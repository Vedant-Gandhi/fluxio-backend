package routes

import "fluxio-backend/pkg/transport/http/controller"

type AuthRoute struct {
	authController *controller.AuthController
}

func NewAuthRoute(authController *controller.AuthController) *AuthRoute {
	return &AuthRoute{
		authController: authController,
	}
}
