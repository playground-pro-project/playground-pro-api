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

func New(db *gorm.DB) user.UserData {
	return userQuery{
		db: db,
	}
}

// Create implements user.UserData.
func (uq userQuery) Create(user user.UserCore) (string, error) {
	// Convert UserCore to UserModel
	userModel := UserCoreToModel(user)

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
	deleteResult := uq.db.Where("user_id = ?", userID).Delete(&User{})
	if deleteResult.Error != nil {
		return fmt.Errorf("failed to delete user: %w", deleteResult.Error)
	}
	if deleteResult.RowsAffected == 0 {
		return fmt.Errorf("no user found with ID: %s", userID)
	}

	return nil
}

// GetByID implements user.UserData.
func (uq userQuery) GetByID(userID string) (user.UserCore, error) {
	var userModel User

	query := uq.db.Preload("Venues").Preload("Reservations").Preload("Reviews").Where("user_id = ?", userID).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return user.UserCore{}, fmt.Errorf("user not found with ID: %s", userID)
		}
		return user.UserCore{}, fmt.Errorf("failed to query user: %w", query.Error)
	}

	userCore := UserModelToCore(userModel)
	return userCore, nil
}

// Login implements user.UserData.
func (uq userQuery) Login(email string, password string) (user.UserCore, string, error) {
	var userModel User

	// Check if the email exists in the database
	query := uq.db.Where("email = ?", email).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return user.UserCore{}, "", errors.New("login failed: invalid email")
		}
		return user.UserCore{}, "", fmt.Errorf("failed to query user: %w", query.Error)
	}

	// Compare the provided password with the stored password
	err := helper.ComparePass([]byte(userModel.Password), []byte(password))
	if err != nil {
		return user.UserCore{}, "", errors.New("login failed: invalid password")
	}

	// Generate an access token
	accessToken, err := middlewares.GenerateToken(userModel.UserID)
	if err != nil {
		return user.UserCore{}, "", fmt.Errorf("failed to generate access token: %w", err)
	}

	userCore := UserModelToCore(userModel)
	return userCore, accessToken, nil
}

// UpdateByID implements user.UserData.
func (uq userQuery) UpdateByID(userID string, updatedUser user.UserCore) error {
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
	userToUpdate := UserCoreToModel(updatedUser)

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
