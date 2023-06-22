package data

import (
	"errors"
	"fmt"
	"strings"

	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

// Create implements user.UserData.
func (uq userQuery) Create(user user.UserEntity) (string, error) {
	// Convert UserEntity to UserModel
	userModel := UserEntityToModel(user)

	// Hash the password
	userModel.Password = helper.HashPass(userModel.Password)

	// Generate a unique user ID
	userModel.UserID = helper.GenerateUserID()

	// Insert the user into the database
	createResult := uq.db.Create(&userModel)
	if createResult.Error != nil {
		// Handle specific errors related to email and phone uniqueness
		switch {
		case strings.Contains(createResult.Error.Error(), "'users.email'"):
			return "", errors.New("email already in use")
		case strings.Contains(createResult.Error.Error(), "'users.phone'"):
			return "", errors.New("phone already in use")
		}
		return "", createResult.Error
	}

	// Check if any row was affected during the insertion
	if createResult.RowsAffected == 0 {
		return "", errors.New("failed to insert, row affected is 0")
	}

	// Return the generated user ID
	return userModel.UserID, nil
}

// DeleteByID implements user.UserData.
func (uq userQuery) DeleteByID(userID string) error {
	deleteResult := uq.db.Delete(&User{}, userID)
	if deleteResult.Error != nil {
		return fmt.Errorf("failed to delete user: %w", deleteResult.Error)
	}
	if deleteResult.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %s", userID)
	}

	return nil
}

// GetAll implements user.UserData.
func (uq userQuery) GetAll() ([]user.UserEntity, error) {
	panic("unimplemented")
}

// GetByID implements user.UserData.
func (uq userQuery) GetByID(userID string) (user.UserEntity, error) {
	var userModel User

	query := uq.db.Preload("Venues").Preload("Reservations").Preload("Reviews").Where("user_id = ?", userID).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return user.UserEntity{}, fmt.Errorf("user not found with ID: %s", userID)
		}
		return user.UserEntity{}, fmt.Errorf("failed to query user: %w", query.Error)
	}

	userEntity := UserModelToEntity(userModel)
	return userEntity, nil
}

// Login implements user.UserData.
func (uq userQuery) Login(email string, password string) (user.UserEntity, string, error) {
	var userModel User

	// Check if the email exists in the database
	query := uq.db.Where("email = ?", email).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return user.UserEntity{}, "", errors.New("login failed: invalid email")
		}
		return user.UserEntity{}, "", fmt.Errorf("failed to query user: %w", query.Error)
	}

	// Compare the provided password with the stored password
	err := helper.ComparePass([]byte(userModel.Password), []byte(password))
	if err != nil {
		return user.UserEntity{}, "", errors.New("login failed: invalid password")
	}

	// Generate an access token
	accessToken, err := middlewares.GenerateToken(userModel.UserID)
	if err != nil {
		return user.UserEntity{}, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	userEntity := UserModelToEntity(userModel)
	return userEntity, accessToken, nil
}

// UpdateByID implements user.UserData.
func (uq userQuery) UpdateByID(userID string, updatedUser user.UserEntity) error {
	var userModel User

	// Retrieve the user from the database
	query := uq.db.Where("user_id = ?", userID).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return fmt.Errorf("user not found with ID: %s", userID)
		}
		return fmt.Errorf("failed to query user: %w", query.Error)
	}

	// Hash the updated password if provided
	if updatedUser.Password != "" {
		updatedUser.Password = helper.HashPass(updatedUser.Password)
	}

	// Convert the updated user entity to a model
	userToUpdate := UserEntityToModel(updatedUser)

	// Perform the update operation
	update := uq.db.Model(&userModel).Updates(userToUpdate)
	if update.Error != nil {
		if strings.Contains(update.Error.Error(), "user not found") {
			return fmt.Errorf("failed to update user: %w", update.Error)
		} else if strings.Contains(update.Error.Error(), "'users.email'") {
			return errors.New("email already in use")
		} else if strings.Contains(update.Error.Error(), "'users.phone'") {
			return errors.New("phone already in use")
		}
	}

	return nil
}

func New(db *gorm.DB) user.UserData {
	return userQuery{
		db: db,
	}
}
