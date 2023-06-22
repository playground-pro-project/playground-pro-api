package service

import (
	"errors"
	"fmt"

	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

type userService struct {
	userData user.UserData
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
		return "", fmt.Errorf("%w", err)
	}

	userID, err := us.userData.Create(user)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return userID, nil
}

// DeleteUserByID implements user.UserService.
func (us *userService) DeleteUserByID(userID string) error {
	err := us.userData.DeleteByID(userID)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

// GetUserByID implements user.UserService.
func (us *userService) GetUserByID(userID string) (user.UserEntity, error) {
	userEntity, err := us.userData.GetByID(userID)
	if err != nil {
		return user.UserEntity{}, fmt.Errorf("error: %w", err)
	}

	return userEntity, nil
}

// Login implements user.UserService.
func (us *userService) Login(email string, password string) (user.UserEntity, string, error) {
	if email == "" {
		return user.UserEntity{}, "", errors.New("email is required")
	}
	if password == "" {
		return user.UserEntity{}, "", errors.New("password is required")
	}

	loggedInUser, accessToken, err := us.userData.Login(email, password)
	if err != nil {
		return user.UserEntity{}, "", fmt.Errorf("%w", err)
	}

	return loggedInUser, accessToken, nil
}

// UpdateUserByID implements user.UserService.
func (us *userService) UpdateUserByID(userID string, updatedUser user.UserEntity) error {
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

func New(repo user.UserData) user.UserService {
	return &userService{
		userData: repo,
	}
}
