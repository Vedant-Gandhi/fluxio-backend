package middleware

import (
	"fluxio-backend/pkg/common/schema"
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
	l            schema.Logger
}

func NewAuthMiddleware(authService *service.UserService, tokenService *service.JWTService, logger schema.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService:  authService,
		tokenService: tokenService,
		l:            logger,
	}
}

func (a *AuthMiddleware) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		userCookie, err := c.Cookie(constants.AuthTokenCookieName)
		logger := a.l

		if err != nil || strings.EqualFold(userCookie, "") {
			logger.Debug("User cookie is not found")
			response.Error(c, response.StatusUnauthorized, "Unauthroized user access.", "User is not authenticated.")
			c.AbortWithError(response.StatusUnauthorized, err)
			return
		}

		logger = logger.With("cookie", userCookie)

		userToken, err := a.tokenService.ValidateToken(userCookie)

		if err != nil || strings.EqualFold(userToken.UserID, "") {
			logger.Debug("User token validation failed", err)
			response.Error(c, response.StatusUnauthorized, "Unauthroized user access.", "User has an invalid token.")
			c.AbortWithError(response.StatusUnauthorized, err)
			return
		}

		userId := model.UserID(userToken.UserID)
		logger = logger.With("user_id", userId)

		user, err := a.authService.GetUserByID(userId)

		// If the user is not found or is blacklisted, return an error.
		if err != nil || user.IsBlackListed {
			logger.Debug("Error occurred when getting user by id", err)
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
