package controller

import (
	"fluxio-backend/pkg/model"
	"fluxio-backend/pkg/service"
	"fluxio-backend/pkg/transport/http/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *service.UserService
}

func NewAuthController(userService *service.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

type rawUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func (a *AuthController) RegisterUser(c *gin.Context) {
	var rawUser rawUserRequest

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

	id, err := a.userService.CreateUser(user, user.Password)

	if err != nil {
		response.Error(c, response.StatusUnprocessableEntity, response.MsgUserCreationFailed, err.Error())
		return
	}

	response.Success(
		c,
		response.StatusCreated,
		"User created sucessfully",
		map[string]any{
			"id": id,
		})

}
