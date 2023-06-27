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

var log = middlewares.Log()

type userHandler struct {
	userService user.UserService
}

func New(service user.UserService) *userHandler {
	return &userHandler{
		userService: service,
	}
}

func (uh *userHandler) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := RegisterRequest{}
		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
		}

		userCore := RegisterRequestToCore(req)
		newUser, err := uh.userService.Register(userCore)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
		}

		userResp := UserCoreToRegisterResponse(newUser)

		return c.JSON(http.StatusCreated, helper.SuccessResponse(userResp, "Check OTP number sent to your email"))
	}
}

func (uh *userHandler) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := LoginRequest{}
		errBind := c.Bind(&request)
		if errBind != nil {
			log.Error("controller - error on bind request")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request"+errBind.Error(), nil, nil))
		}

		resp, token, err := uh.userService.Login(RequestToCore(request))
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "invalid email format"):
				log.Error("bad request, invalid email format")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request, invalid email format", nil, nil))
			case strings.Contains(err.Error(), "password cannot be empty"):
				log.Error("bad request, password cannot be empty")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request, password cannot be empty", nil, nil))
			case strings.Contains(err.Error(), "invalid email and password"):
				log.Error("bad request, invalid email and password")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request, invalid email and password", nil, nil))
			case strings.Contains(err.Error(), "password does not match"):
				log.Error("bad request, password does not match")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request, password does not match", nil, nil))
			case strings.Contains(err.Error(), "error while creating jwt token"):
				log.Error("internal server error, error while creating jwt token")
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
			default:
				log.Error("internal server error")
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
			}
		}

		loginResp := UserCoreToLoginResponse(resp)
		loginResp.Token = token

		if loginResp.AccountStatus == "unverified" {
			loginResp.Token = ""
			return c.JSON(http.StatusOK, helper.SuccessResponse(loginResp, "OTP validation is required"))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse(loginResp, "Login success"))
	}
}

func (uh *userHandler) ReSendOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := LoginRequest{}
		errBind := c.Bind(&req)
		if errBind != nil {
			log.Error("controller - error on bind request")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request"+errBind.Error(), nil, nil))
		}

		userID, _ := uh.userService.GetUserID(req.Email)
		user, err := uh.userService.GetByID(userID)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error(err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			}
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
		}

		err = uh.userService.StoreToRedis(user)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
		}

		return c.JSON(http.StatusCreated, helper.SuccessResponse(nil, "Check OTP number sent to your email"))
	}
}

func (uh *userHandler) ValidateOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := OTPInputReq{}
		err := c.Bind(&req)
		if err != nil {
			log.Error("error binding request")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": "error binding request",
			})
		}

		isValid, err := uh.userService.VerifyOTP(req.UserID, req.OTP)
		if !isValid {
			log.Error("OTP has been expired")
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"message": "OTP has been expired",
			})
		}

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}

		err = uh.userService.UpdateByID(req.UserID, user.UserCore{
			AccountStatus: "verified",
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error, please try again later"))
		}

		user, err := uh.userService.GetByID(req.UserID)

		resp, token, _ := uh.userService.Login(user)

		loginResp := UserCoreToLoginResponse(resp)
		loginResp.Token = token

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    loginResp,
			"message": "Verification success",
		})
	}
}

func (uh *userHandler) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := middlewares.ExtractToken(c)
		if err != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		user, err := uh.userService.GetByID(userID)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			}
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
		}

		userResponse := UserCoreToGetUserResponse(user)

		return c.JSON(http.StatusOK, helper.SuccessResponse(userResponse, "User profile retrieved successfully"))
	}
}

func (uh *userHandler) UpdatePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
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
}

func (uh *userHandler) UpdateUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
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
}
func (uh *userHandler) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
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
}

func (uh *userHandler) UploadProfilePicture() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		awsService := aws.InitS3()

		file, err := c.FormFile("profile_picture")
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
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
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
		}
		defer fileContent.Close()

		err = awsService.UploadFile(path, fileType, fileContent)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
		}

		var updatedUser user.UserCore
		updatedUser.ProfilePicture = fmt.Sprintf("%s%s", profilePictureBaseURL, filepath.Base(file.Filename))

		err = uh.userService.UpdateByID(userId, updatedUser)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error(err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			}
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Profile picture updated successfully"))
	}
}

func (uh *userHandler) RemoveProfilePicture() echo.HandlerFunc {
	return func(c echo.Context) error {
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
				log.Error(err.Error())
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"error": err.Error(),
				})
			}
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Profile picture removed successfully",
		})
	}
}

func (uh *userHandler) UploadOwnerFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		awsService := aws.InitS3()

		file, err := c.FormFile("owner_docs")
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
		}

		// Check file size before opening it
		fileSize := file.Size
		if fileSize > maxOwnerFileSize {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Please upload file smaller than 5 MB."))
		}

		// Get the file type from the Content-Type header
		fileType := file.Header.Get("Content-Type")

		path := "owner-docs/" + file.Filename
		fileContent, err := file.Open()
		if err != nil {
			log.Error(err.Error())
			return err
		}
		defer fileContent.Close()

		err = awsService.UploadFile(path, fileType, fileContent)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
		}

		var updatedUser user.UserCore
		updatedUser.OwnerFile = fmt.Sprintf("%s%s", ownerFileBaseURL, filepath.Base(file.Filename))
		updatedUser.Role = "owner"

		err = uh.userService.UpdateByID(userId, updatedUser)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error(err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			}
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "File added successfully"))
	}
}
