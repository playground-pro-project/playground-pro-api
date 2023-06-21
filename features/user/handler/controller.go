package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/features/user"
)

type userHandler struct {
	userService user.UserService
}

func New(service user.UserService) *userHandler {
	return &userHandler{
		userService: service,
	}
}

func (uh *userHandler) Register(c echo.Context) error {
	req := RegisterRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}

	userEntity := UserRequestToEntity(req)
	_, err = uh.userService.CreateUser(userEntity)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User registered successfully",
	})
}
