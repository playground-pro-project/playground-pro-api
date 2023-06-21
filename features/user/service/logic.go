package service

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

type userService struct {
	userData user.UserData
	validate *validator.Validate
}

// CreateUser implements user.UserService.
func (us *userService) CreateUser(user user.UserEntity) (string, error) {
	if user.Fullname == "" {
		return "", errors.New("error, name is required")
	}
	if user.Email == "" {
		return "", errors.New("error, mail is required")
	} else if _, err := helper.ValidateMailAddress(user.Email); !err {
		return "", errors.New("error, invalid email format")
	}
	if user.Phone == "" {
		return "", errors.New("error, phone is required")
	}
	if user.Password == "" {
		return "", errors.New("error, password is required")
	}

	err := helper.ValidatePassword(user.Password)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	userID, err := us.userData.Create(user)
	if err != nil {
		return "", fmt.Errorf("%v", err)
	}

	return userID, nil
}

// DeleteUserByID implements user.UserService.
func (us *userService) DeleteUserByID(userID string) error {
	panic("unimplemented")
}

// GetAllUsers implements user.UserService.
func (us *userService) GetAllUsers() ([]user.UserEntity, error) {
	panic("unimplemented")
}

// GetUserByID implements user.UserService.
func (us *userService) GetUserByID(userID string) (user.UserEntity, error) {
	panic("unimplemented")
}

// Login implements user.UserService.
func (us *userService) Login(email string, password string) (user.UserEntity, string, error) {
	panic("unimplemented")
}

// UpdateUserByID implements user.UserService.
func (us *userService) UpdateUserByID(userID string, updatedUser user.UserEntity) error {
	panic("unimplemented")
}

func New(repo user.UserData) user.UserService {
	return &userService{
		userData: repo,
		validate: validator.New(),
	}
}
