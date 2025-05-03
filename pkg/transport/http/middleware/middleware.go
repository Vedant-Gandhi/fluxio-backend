package middleware

import (
	"fluxio-backend/pkg/constants"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService  service.UserService
	tokenService service.JWTService
}

func (a *AuthMiddleware) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCookie, err := c.Cookie(constants.AuthTokenCookieName)

		if err != nil || strings.EqualFold(userCookie, "") {
			response.Error(c, response.StatusUnauthorized, "Unauthroized user access.", "User is not authenticated.")
			c.AbortWithError(response.StatusUnauthorized, err)
			return
		}

		userToken, err := a.tokenService.ValidateToken(userCookie)

		if err != nil || strings.EqualFold(userToken.UserID, "") {
			response.Error(c, response.StatusUnauthorized, "Unauthroized user access.", "User has an invalid token.")
			c.AbortWithError(response.StatusUnauthorized, err)
			return
		}

		userId := model.UserID(userToken.UserID)

	}
}
