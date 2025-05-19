package controller

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

type AuthController struct {
	userService *service.UserService
	l           schema.Logger
}

func NewAuthController(userService *service.UserService, logger schema.Logger) *AuthController {
	return &AuthController{
		userService: userService,
		l:           logger,
	}
}

type userRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (a *AuthController) RegisterUser(c *gin.Context) {
	logger := a.l.WithField(map[string]interface{}{
		"client_ip": c.ClientIP(),
		"path":      c.FullPath(),
	})

	var rawUser userRegisterRequest

	if err := c.ShouldBindJSON(&rawUser); err != nil {
		logger.Warn("Invalid registration payload")
		response.Error(c, response.StatusBadRequest, "Invalid request payload", "The payload is not valid.")
		return
	}

	logger = logger.With("username", rawUser.Username).With("email", rawUser.Email)

	if strings.EqualFold(rawUser.Password, "") {
		logger.Warn("Empty password in registration")
		response.Error(c, response.StatusBadRequest, "Invalid Password", "The password sent by the user is empty.")
		return
	}

	user := model.User{
		Username: rawUser.Username,
		Password: rawUser.Password,
		Email:    rawUser.Email,
	}

	id, token, err := a.userService.CreateUser(user, user.Password)

	if err != nil {
		if err == fluxerrors.ErrUsernameExists {
			logger.Info("Registration failed - username already exists")
			response.Error(c, response.StatusConflict, "User already exists", err.Error())
			return
		} else if err == fluxerrors.ErrEmailExists {
			logger.Info("Registration failed - email already exists")
			response.Error(c, response.StatusConflict, "User already exists", err.Error())
			return
		}
		logger.Error("User creation failed", err)
		response.Error(c, response.StatusUnprocessableEntity, response.MsgUserCreationFailed, err.Error())
		return
	}

	c.SetCookie(constants.AuthTokenCookieName, token, constants.AuthTokenCookieExp, "/", "", false, true)

	logger.Info("User registered successfully", "user_id", id.String())
	response.Success(
		c,
		response.StatusCreated,
		"User created sucessfully",
		map[string]any{
			"id": id,
		})
}

type userLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

func (a *AuthController) LoginUser(c *gin.Context) {
	logger := a.l.WithField(map[string]interface{}{
		"client_ip": c.ClientIP(),
		"path":      c.FullPath(),
	})

	var loginData userLoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		logger.Warn("Invalid login payload")
		response.Error(c, response.StatusBadRequest, "Invalid request payload", "The payload is not valid.")
		return
	}

	// Add user identifier to logger context
	if loginData.Email != "" {
		logger = logger.With("email", loginData.Email)
	} else if loginData.Username != "" {
		logger = logger.With("username", loginData.Username)
	}

	// Check for password empty and username or email empty.
	if strings.EqualFold(loginData.Password, "") || (strings.EqualFold(loginData.Username, "") && strings.EqualFold(loginData.Email, "")) {
		logger.Warn("Login attempt with empty credentials")
		response.Error(c, response.StatusBadRequest, "Invalid Credentials", "The user credentials are not valid.")
		return
	}

	user := model.User{
		Username: loginData.Username,
		Password: loginData.Password,
		Email:    loginData.Email,
	}

	res, token, err := a.userService.Login(user)

	if err != nil {
		if err == fluxerrors.ErrInvalidCredentials {
			logger.Info("Login failed - invalid credentials")
			response.Error(c, response.StatusUnauthorized, "Invalid credentials", "The user credentials are not valid.")
			return
		}
		if err == fluxerrors.ErrUserNotFound {
			logger.Info("Login failed - user not found")
			response.Error(c, response.StatusNotFound, "User not found", "The user does not exist.")
			return
		}

		logger.Error("Login failed - server error", err)
		response.Error(c, response.StatusInternalServerError, "Internal server error", err.Error())
		return
	}

	// Set token
	c.SetCookie(constants.AuthTokenCookieName, token, constants.AuthTokenCookieExp, "/", "", false, true)

	logger.Info("User logged in successfully", "user_id", res.ID.String())
	response.Success(
		c,
		response.StatusOK,
		"User logged in successfully",
		res,
	)
}
