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

const (
	maxFileSize              = 1 << 20 // 1 MB
	profilePictureBaseURL    = "https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/user-profile-picture/"
	defaultProfilePictureURL = "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
	ownerFileBaseURL         = "https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/owner-docs/"
)

var log = middlewares.Log()

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
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
	}

	userCore := RegisterRequestToCore(req)
	_, err = uh.userService.CreateUser(userCore)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "User registered successfully"))
}

func (uh *userHandler) Login(c echo.Context) error {
	req := LoginRequest{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	user, token, err := uh.userService.Login(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	response := map[string]interface{}{
		"user_id": user.UserID,
		"token":   token,
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(response, "Login successful"))
}

func (uh *userHandler) GetUserProfile(c echo.Context) error {
	userId, err := middlewares.ExtractToken(c)
	if err != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	user, err := uh.userService.GetByID(userId)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
	}

	userResponse := UserCoreToGetUserResponse(user)

	return c.JSON(http.StatusOK, helper.SuccessResponse(userResponse, "User profile retrieved successfully"))
}

func (uh *userHandler) UpdatePassword(c echo.Context) error {
	req := ChangePasswordRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
	}

	userId, err := middlewares.ExtractToken(c)
	if err != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	user, err := uh.userService.GetByID(userId)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error, please try again later"))
	}

	err = helper.ComparePass([]byte(user.Password), []byte(req.OldPassword))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Wrong password"))
	}

	// req.NewPassword = helper.HashPass(req.NewPassword)

	updatedUser := UpdatePasswordRequestToCore(req)
	err = uh.userService.UpdateByID(userId, updatedUser)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error, please try again later"))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Password updated successfully"))
}

func (uh *userHandler) UpdateUserProfile(c echo.Context) error {
	req := EditProfileRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
	}

	userId, err := middlewares.ExtractToken(c)
	if err != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}
	updatedUser := EditProfileRequestToCore(req)

	err = uh.userService.UpdateByID(userId, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(
				http.StatusNotFound, helper.ErrorResponse(err.Error()),
			)
		} else if strings.Contains(err.Error(), "email") {
			return c.JSON(
				http.StatusBadRequest, helper.ErrorResponse(err.Error()),
			)
		} else if strings.Contains(err.Error(), "phone") {
			return c.JSON(
				http.StatusBadRequest, helper.ErrorResponse(err.Error()),
			)
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "User profile updated successfully"))
}

func (uh *userHandler) DeleteUser(c echo.Context) error {
	userId, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	err := uh.userService.DeleteByID(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusBadRequest, helper.SuccessResponse(nil, "User account deleted successfully"))
}

func (uh *userHandler) UploadProfilePicture(c echo.Context) error {
	userId, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	awsService := aws.InitS3()

	file, err := c.FormFile("profile_picture")
	if err != nil {
		return err
	}

	// Check file size before opening it
	fileSize := file.Size
	if fileSize > maxFileSize {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Please upload a picture smaller than 1 MB."))
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
	updatedUser.ProfilePicture = fmt.Sprintf("%s%s", profilePictureBaseURL, filepath.Base(file.Filename))

	err = uh.userService.UpdateByID(userId, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Profile picture updated successfully"))
}

func (uh *userHandler) RemoveProfilePicture(c echo.Context) error {
	userId, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	updatedUser := user.UserCore{
		ProfilePicture: defaultProfilePictureURL,
	}
	err := uh.userService.UpdateByID(userId, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Profile picture removed successfully"))
}

func (uh *userHandler) UploadOwnerFile(c echo.Context) error {
	userId, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	awsService := aws.InitS3()

	file, err := c.FormFile("owner_docs")
	if err != nil {
		return err
	}

	// Check file size before opening it
	fileSize := file.Size
	if fileSize > maxFileSize {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Please upload file smaller than 1 MB."))
	}

	// Get the file type from the Content-Type header
	fileType := file.Header.Get("Content-Type")

	path := "owner-docs/" + file.Filename
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
	updatedUser.OwnerFile = fmt.Sprintf("%s%s", ownerFileBaseURL, filepath.Base(file.Filename))

	// updatedUser.Role = "owner"

	err = uh.userService.UpdateByID(userId, updatedUser)
	if err != nil {
		if strings.Contains(err.Error(), "user not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "File added successfully"))
}
