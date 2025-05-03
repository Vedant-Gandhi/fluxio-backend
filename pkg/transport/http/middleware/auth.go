package middleware

import (
	"fluxio-backend/pkg/constants"
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService  *service.UserService
	tokenService *service.JWTService
}

func NewAuthMiddleware(authService *service.UserService, tokenService *service.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		authService:  authService,
		tokenService: tokenService,
	}
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

		user, err := a.authService.GetUserByID(userId)

		// If the user is not found or is blacklisted, return an error.
		if err != nil || user.IsBlackListed {
			if err == fluxerrors.ErrUserNotFound || err == fluxerrors.ErrInvalidUserID {
				response.Error(c, response.StatusNotFound, "User not found.", "User not found.")
				c.AbortWithError(response.StatusUnauthorized, err)
				return
			}

			response.Error(c, response.StatusUnauthorized, "Unauthroized user access.", "User is not authenticated.")
			c.AbortWithError(response.StatusUnauthorized, err)
			return
		}

		c.Set(constants.GinUserContextKey, user)
		c.Next()

	}
}
