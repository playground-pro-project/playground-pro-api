package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/utils/aws"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
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

	userCore := RegisterRequestToCore(req)
	_, err = uh.userService.CreateUser(userCore)
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

	user, err := uh.userService.GetByID(userID)
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

	userResponse := UserCoreToGetUserResponse(user)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    userResponse,
		"message": "User profile retrieved successfully",
	})
}

func (uh *userHandler) UpdatePassword(c echo.Context) error {
	req := ChangePasswordRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}

	userID := middlewares.ExtractUserIDFromToken(c)
	user, err := uh.userService.GetByID(userID)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error, please try again later",
		})
	}

	err = helper.ComparePass([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Wrong password",
		})
	}

	// req.NewPassword = helper.HashPass(req.NewPassword)

	updatedUser := UpdatePasswordRequestToCore(req)
	err = uh.userService.UpdateByID(userID, updatedUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error, please try again later",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Password updated successfully",
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
	updatedUser := EditProfileRequestToCore(req)

	err = uh.userService.UpdateByID(userID, updatedUser)
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

func (uh *userHandler) DeleteUser(c echo.Context) error {
	userID := middlewares.ExtractUserIDFromToken(c)
	err := uh.userService.DeleteByID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": "User account deleted successfully",
	})
}

const (
	MaxFileSize = 1 << 20 // 1 MB
)

func (uh *userHandler) UploadProfilePicture(c echo.Context) error {
	userID := middlewares.ExtractUserIDFromToken(c)

	awsService := aws.InitS3()

	file, err := c.FormFile("profile_picture")
	if err != nil {
		return err
	}

	// Check file size before opening it
	fileSize := file.Size
	if fileSize > MaxFileSize {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "Please upload a picture smaller than 1 MB.",
		})
	}

	// Get the file type from the Content-Type header
	fileType := file.Header.Get("Content-Type")

	path := "profile-picture/" + file.Filename
	fileContent, err := file.Open()
	if err != nil {
		return err
	}
	defer fileContent.Close()

	err = awsService.UploadFile(path, fileType, fileContent)
	if err != nil {
		return err
	}

	var updatedUser user.UserCore
	updatedUser.ProfilePicture = fmt.Sprintf(
		"https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/profile-picture/%s",
		filepath.Base(file.Filename),
	)

	err = uh.userService.UpdateByID(userID, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Profile picture updated successfully",
	})
}

func (uh *userHandler) RemoveProfilePicture(c echo.Context) error {
	userID := middlewares.ExtractUserIDFromToken(c)

	updatedUser := user.UserCore{
		ProfilePicture: "https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/profile-picture/default-image.jpg",
	}
	err := uh.userService.UpdateByID(userID, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Profile picture removed successfully",
	})
}

func (uh *userHandler) UploadOwnerFile(c echo.Context) error {
	panic("unimplemented")
}
