package controller

import (
	fluxerrors "fluxio-backend/pkg/errors"
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	TOKEN_COOKIE_NAME = "token"
	TOKEN_COOKIE_EXP  = 8 * 3600
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

type userRegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (a *AuthController) RegisterUser(c *gin.Context) {
	var rawUser userRegisterRequest

	if err := c.ShouldBindJSON(&rawUser); err != nil {
		response.Error(c, response.StatusBadRequest, "Invalid request payload", "The payload is not valid.")
		return
	}

	if strings.EqualFold(rawUser.Password, "") {
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
			response.Error(c, response.StatusConflict, "User already exists", err.Error())
			return
		}
		response.Error(c, response.StatusUnprocessableEntity, response.MsgUserCreationFailed, err.Error())
		return
	}

	c.SetCookie(TOKEN_COOKIE_NAME, token, TOKEN_COOKIE_EXP, "/", "", false, true)

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
	var loginData userLoginRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		response.Error(c, response.StatusBadRequest, "Invalid request payload", "The payload is not valid.")
		return
	}

	// Check for password empty and username or email empty.
	if strings.EqualFold(loginData.Password, "") || (strings.EqualFold(loginData.Username, "") && strings.EqualFold(loginData.Email, "")) {
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
			response.Error(c, response.StatusUnauthorized, "Invalid credentials", "The user credentials are not valid.")
			return
		}
		if err == fluxerrors.ErrUserNotFound {
			response.Error(c, response.StatusNotFound, "User not found", "The user does not exist.")
			return
		}

		response.Error(c, response.StatusInternalServerError, "Internal server error", err.Error())
	}

	// Set token
	c.SetCookie(TOKEN_COOKIE_NAME, token, TOKEN_COOKIE_EXP, "/", "", false, true)

	response.Success(
		c,
		response.StatusOK,
		"User logged in successfully",
		res,
	)
}
