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
		newUser, otp, err := uh.userService.Register(userCore)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
		}

		userResp := UserCoreToRegisterResponse(newUser)
		userResp.OTP = otp

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
				return helper.BadRequestError(c, "Bad request, invalid email format")
			case strings.Contains(err.Error(), "password cannot be empty"):
				log.Error("bad request, password cannot be empty")
				return helper.BadRequestError(c, "Bad request, password cannot be empty")
			case strings.Contains(err.Error(), "invalid email and password"):
				log.Error("bad request, invalid email and password")
				return helper.BadRequestError(c, "Bad request, invalid email and password")
			case strings.Contains(err.Error(), "password does not match"):
				log.Error("bad request, password does not match")
				return helper.BadRequestError(c, "Bad request, password does not match")
			case strings.Contains(err.Error(), "no row affected"):
				log.Error("no row affected")
				return helper.NotFoundError(c, "The requested resource was not found")
			case strings.Contains(err.Error(), "error while creating jwt token"):
				log.Error("internal server error, error while creating jwt token")
				return helper.InternalServerError(c, "Internal server error")
			default:
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
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
		req := ResendOtpReq{}
		err := c.Bind(&req)
		if err != nil {
			log.Error("controller - error on bind request")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("controller - error on bind request"+err.Error()))
		}

		userID, _ := uh.userService.GetUserID(req.Email)
		user, err := uh.userService.GetByID(userID)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error(err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			}
			log.Error(err.Error())
			return helper.InternalServerError(c, "Internal server error")
		}

		err = uh.userService.StoreToRedis(user)
		if err != nil {
			log.Error(err.Error())
			return helper.InternalServerError(c, "Internal server error")
		}

		return c.JSON(http.StatusCreated, helper.SuccessResponse(nil, "Check OTP number sent to your email"))
	}
}

func (uh *userHandler) ValidateOTP() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := VerifyReq{}
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
				"status":  "Fail",
				"message": "OTP has been expired",
			})
		}

		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"status": "Fail",
				"error":  err.Error(),
			})
		}

		err = uh.userService.UpdateByID(req.UserID, user.UserCore{
			AccountStatus: "verified",
		})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error, please try again later"))
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"status":  "Success",
			"message": "Verification success",
		})
	}
}

func (uh *userHandler) GetUserProfile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
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
}

func (uh *userHandler) UpdatePassword() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := ChangePasswordRequest{}

		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		err := c.Bind(&req)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
		}

		if req.OldPassword == "" {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Old password is required"))
		}

		if req.NewPassword == "" {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("New password is required"))
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

		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		// Check if the request body is empty
		if c.Request().ContentLength == 0 {
			log.Error("Empty request payload")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Empty request payload"))
		}

		err := c.Bind(&req)
		if err != nil {
			log.Error("Invalid request payload")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
		}

		updatedUser := EditProfileRequestToCore(req)
		fmt.Println("Updated User: ", updatedUser)

		err = uh.userService.UpdateByID(userId, updatedUser)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			} else if strings.Contains(err.Error(), "email") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
			} else if strings.Contains(err.Error(), "phone") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
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
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		err := uh.userService.DeleteByID(userId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "User account deleted successfully"))
	}
}

func (uh *userHandler) UploadProfilePicture() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		awsService := aws.InitS3()

		file, err := c.FormFile("profile_picture")
		if err != nil {
			log.Error("Failed to retrieve profile picture: " + err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Failed to retrieve profile picture: "+err.Error()))
		}

		// Get the file type from the Content-Type header
		fileType := file.Header.Get("Content-Type")
		fileExtension := filepath.Ext(file.Filename)
		fileExtension = strings.ToLower(fileExtension)

		allowedExtensions := []string{".jpg", ".jpeg", ".png"}
		extensionAllowed := false
		for _, ext := range allowedExtensions {
			if ext == fileExtension {
				extensionAllowed = true
				break
			}
		}
		if !extensionAllowed {
			log.Error(fileExtension + " is not allowed")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(fileExtension+" is not allowed. Only JPG, JPEG, and PNG files are allowed."))
		}

		// Check file size before opening it
		fileSize := file.Size
		if fileSize > maxFileSize {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Please upload a picture smaller than 1 MB."))
		}

		id := helper.GenerateIdentifier()
		filename := id + "-" + file.Filename
		path := "profile-picture/" + filename

		fileContent, err := file.Open()
		if err != nil {
			log.Error("Failed to open file: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to open file: "+err.Error()))
		}
		defer fileContent.Close()

		// Upload profile picture file to cloud
		err = awsService.UploadFile(path, fileType, fileContent)
		if err != nil {
			log.Error("Failed to upload file to cloud service: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to upload file to cloud service: "+err.Error()))
		}

		usr, err := uh.userService.GetByID(userId)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("User not found: "+err.Error()))
		}

		// Delete profile picture in the cloud before updated
		prevFilename := filepath.Base(usr.ProfilePicture)
		prevPath := "profile-picture/" + prevFilename

		err = awsService.DeleteFile(prevPath)
		if err != nil {
			log.Error("Failed to delete file from cloud service: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to delete file from cloud service: "+err.Error()))
		}

		// Update user profile picture in database
		var updatedUser user.UserCore
		updatedUser.ProfilePicture = fmt.Sprintf("%s%s", profilePictureBaseURL, filepath.Base(filename))

		err = uh.userService.UpdateByID(userId, updatedUser)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error("User not found: " + err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("User not found: "+err.Error()))
			}
			log.Error("Failed to update user: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to update user: "+err.Error()))
		}

		resp := map[string]interface{}{
			"profile_picture": updatedUser.ProfilePicture,
		}

		log.Sugar().Infof(userId + " updated profile picture successfully")
		return c.JSON(http.StatusOK, helper.SuccessResponse(resp, "Profile picture updated successfully"))
	}
}

func (uh *userHandler) RemoveProfilePicture() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		awsService := aws.InitS3()

		usr, err := uh.userService.GetByID(userId)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("User not found: "+err.Error()))
		}

		// Delete profile picture in the cloud before updated
		prevFilename := filepath.Base(usr.ProfilePicture)
		prevPath := "profile-picture/" + prevFilename

		err = awsService.DeleteFile(prevPath)
		if err != nil {
			log.Error("Failed to delete file from cloud service: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to delete file from cloud service: "+err.Error()))
		}

		updatedUser := user.UserCore{
			ProfilePicture: defaultProfilePictureURL,
		}
		err = uh.userService.UpdateByID(userId, updatedUser)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error(err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
			}
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
		}

		log.Sugar().Infof(userId + " removed profile picture successfully")
		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Profile picture removed successfully"))
	}
}

func (uh *userHandler) UploadOwnerFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		awsService := aws.InitS3()

		file, err := c.FormFile("owner_docs")
		if err != nil {
			log.Error("Failed to retrieve file: " + err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Failed to retrieve file: "+err.Error()))
		}

		// Get the file type from the Content-Type header
		fileType := file.Header.Get("Content-Type")
		fileExtension := filepath.Ext(file.Filename)
		fileExtension = strings.ToLower(fileExtension)

		allowedExtensions := []string{".jpg", ".jpeg", ".png", ".pdf"}
		extensionAllowed := false
		for _, ext := range allowedExtensions {
			if ext == fileExtension {
				extensionAllowed = true
				break
			}
		}
		if !extensionAllowed {
			log.Error(fileExtension + " is not allowed")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(fileExtension+" is not allowed. Only JPG, JPEG, PNG and PDF files are allowed."))
		}

		// Check file size before opening it
		fileSize := file.Size
		if fileSize > maxOwnerFileSize {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Please upload a file smaller than 5 MB."))
		}

		id := helper.GenerateIdentifier()
		filename := id + "-" + file.Filename
		path := "owner-docs/" + filename

		fileContent, err := file.Open()
		if err != nil {
			log.Error("Failed to open file: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to open file: "+err.Error()))
		}
		defer fileContent.Close()

		err = awsService.UploadFile(path, fileType, fileContent)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
		}

		usr, err := uh.userService.GetByID(userId)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, helper.ErrorResponse("User not found: "+err.Error()))
		}

		// Delete profile picture in the cloud before updated
		prevFilename := filepath.Base(usr.ProfilePicture)
		prevPath := "owner-docs/" + prevFilename

		err = awsService.DeleteFile(prevPath)
		if err != nil {
			log.Error("Failed to delete file from cloud service: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to delete file from cloud service: "+err.Error()))
		}

		var updatedUser user.UserCore
		updatedUser.OwnerFile = fmt.Sprintf("%s%s", ownerFileBaseURL, filepath.Base(filename))
		updatedUser.Role = "owner"

		err = uh.userService.UpdateByID(userId, updatedUser)
		if err != nil {
			if strings.Contains(err.Error(), "user not found") {
				log.Error("User not found: " + err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("User not found: "+err.Error()))
			}
			log.Error("Failed to update user: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to update user: "+err.Error()))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "File added successfully"))
	}
}
