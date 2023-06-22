package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
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

	userEntity := RegisterRequestToEntity(req)
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

func (uh *userHandler) Login(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	user, token, err := uh.userService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	response := map[string]interface{}{
		"user_id": user.UserID,
		"token":   token,
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    response,
		"message": "Login successful",
	})
}

func (uh *userHandler) GetUserProfile(c echo.Context) error {
	userID := middlewares.ExtractUserIDFromToken(c)

	user, err := uh.userService.GetUserByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	userResponse := UserEntityToGetUserResponse(user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    userResponse,
		"message": "User profile retrieved successfully",
	})
}

func (uh *userHandler) UpdateUserProfile(c echo.Context) error {
	req := EditProfileRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}

	userID := middlewares.ExtractUserIDFromToken(c)
	updatedUser := EditProfileRequestToEntity(req)

	err = uh.userService.UpdateUserByID(userID, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		} else if strings.Contains(err.Error(), "email") {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		} else if strings.Contains(err.Error(), "phone") {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User profile updated successfully",
	})
}
