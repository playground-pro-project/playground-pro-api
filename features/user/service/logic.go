package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

var log = middlewares.Log()

type userService struct {
	userData user.UserData
	validate *validator.Validate
}

func New(ud user.UserData, v *validator.Validate) user.UserService {
	return &userService{
		userData: ud,
		validate: v,
	}
}

// Register implements user.UserService
func (us *userService) Register(request user.UserCore) (user.UserCore, error) {
	err := us.validate.Struct(request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Fullname"):
			log.Warn("fullname cannot be empty")
			return user.UserCore{}, errors.New("fullname cannot be empty")
		case strings.Contains(err.Error(), "Email"):
			log.Warn("invalid email format")
			return user.UserCore{}, errors.New("invalid email format")
		case strings.Contains(err.Error(), "Phone"):
			log.Warn("phone cannot be empty")
			return user.UserCore{}, errors.New("phone cannot be empty")
		case strings.Contains(err.Error(), "Password"):
			log.Warn("password cannot be empty")
			return user.UserCore{}, errors.New("password cannot be empty")
		}
	}

	result, err := us.userData.Register(request)
	if err != nil {
		message := ""
		switch {
		case strings.Contains(err.Error(), "error while hashing password"):
			log.Error("error while hashing password")
			message = "error while hashing password"
		case strings.Contains(err.Error(), "error insert data, duplicated"):
			log.Error("error insert data, duplicated")
			message = "error insert data, duplicated"
		default:
			log.Error("internal server error")
			message = "internal server error"
		}
		log.Error("request cannot be empty")
		return user.UserCore{}, errors.New(message)
	}

	log.Sugar().Infof("new user has been created: %s", result.Email)
	return result, nil
}

// Login implements user.UserService
func (us *userService) Login(request user.UserCore) (user.UserCore, string, error) {
	err := us.validate.Struct(request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Email"):
			log.Warn("invalid email format")
			return user.UserCore{}, "", errors.New("invalid email format")
		case strings.Contains(err.Error(), "Password"):
			log.Warn("password cannot be empty")
			return user.UserCore{}, "", errors.New("password cannot be empty")
		}
	}

	result, token, err := us.userData.Login(request)
	if err != nil {
		message := ""
		switch {
		case strings.Contains(err.Error(), "invalid email and password"):
			log.Error("invalid email and password")
			message = "invalid email and password"
		case strings.Contains(err.Error(), "password does not match"):
			log.Error("password does not match")
			message = "password does not match"
		case strings.Contains(err.Error(), "error while creating jwt token"):
			log.Error("error while creating jwt token")
			message = "error while creating jwt token"
		default:
			log.Error("internal server error")
			message = "internal server error"
		}
		return user.UserCore{}, "", errors.New(message)
	}

	log.Sugar().Infof("user has been logged in: %s", result.UserID)
	return result, token, nil
}

// GenerateOTP implements user.UserService.
func (us *userService) GenerateOTP(request user.UserCore) (user.UserCore, error) {
	result, err := us.userData.GenerateOTP(request)
	if err != nil {
		return user.UserCore{}, errors.New("user not found")
	}

	return result, err
}

// VerifyOTP implements user.UserService.
func (us *userService) VerifyOTP(request user.UserCore) (user.UserCore, error) {
	result, err := us.userData.VerifyOTP(request)
	if err != nil {
		return user.UserCore{}, errors.New("user not found")
	}

	return result, err
}

// DeleteUserByID implements user.UserService.
func (us *userService) DeleteByID(userID string) error {
	err := us.userData.DeleteByID(userID)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

// GetUserByID implements user.UserService.
func (us *userService) GetByID(userID string) (user.UserCore, error) {
	userEntity, err := us.userData.GetByID(userID)
	if err != nil {
		return user.UserCore{}, fmt.Errorf("error: %w", err)
	}

	return userEntity, nil
}

// UpdateUserByID implements user.UserService.
func (us *userService) UpdateByID(userID string, updatedUser user.UserCore) error {
	if updatedUser.Password != "" {
		err := helper.ValidatePassword(updatedUser.Password)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}
	if updatedUser.Email != "" {
		_, err := helper.ValidateMailAddress(updatedUser.Email)
		if !err {
			return errors.New("error: invalid email format")
		}
	}

	err := us.userData.UpdateByID(userID, updatedUser)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}
