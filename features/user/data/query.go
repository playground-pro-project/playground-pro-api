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

var log = middlewares.Log()

type userQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) user.UserData {
	return &userQuery{
		db: db,
	}
}

// Register implements user.UserData
func (uq *userQuery) Register(request user.UserCore) (user.UserCore, error) {
	hashedPass, err := helper.HashPassword(request.Password)
	if err != nil {
		log.Error("error while hashing password")
		return user.UserCore{}, errors.New("error while hashing password")
	}

	request.Password = hashedPass
	req := UserCoreToModel(request)
	query := uq.db.Table("users").Create(&req)
	if query.Error != nil {
		log.Error("error insert data, duplicated")
		return user.UserCore{}, errors.New("error insert data, duplicated")
	}

	rowAffect := query.RowsAffected
	if rowAffect == 0 {
		log.Warn("no user has been created")
		return user.UserCore{}, errors.New("no row affected")
	}

	log.Sugar().Infof("new user has been created: %s", req.Email)
	return UserModelToCore(req), nil
}

// Login implements user.UserData
func (uq *userQuery) Login(request user.UserCore) (user.UserCore, string, error) {
	result := User{}
	query := uq.db.Table("users").Where("email = ?", request.Email).First(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("user record not found, invalid email and password")
		return user.UserCore{}, "", errors.New("invalid email and password")
	}

	rowAffect := query.RowsAffected
	if rowAffect == 0 {
		log.Warn("no user has been created")
		return user.UserCore{}, "", errors.New("no row affected")
	}

	if !helper.MatchPassword(request.Password, result.Password) {
		log.Warn("password does not match")
		return user.UserCore{}, "", errors.New("password does not match")
	}

	token, err := middlewares.GenerateToken(result.UserID)
	if err != nil {
		log.Error("error while creating jwt token")
		return user.UserCore{}, "", errors.New("error while creating jwt token")
	}

	log.Sugar().Infof("user has been logged in: %s", result.UserID)
	return UserModelToCore(result), token, nil
}

// DeleteByID implements user.UserData.
func (uq *userQuery) DeleteByID(userID string) error {
	deleteResult := uq.db.Where("user_id = ?", userID).Delete(&User{})
	if deleteResult.Error != nil {
		log.Sugar().Errorf("failed to delete user: %w", deleteResult.Error)
		return fmt.Errorf("failed to delete user: %w", deleteResult.Error)
	}
	if deleteResult.RowsAffected == 0 {
		log.Sugar().Errorf("no user found with ID: %s", userID)
		return fmt.Errorf("no user found with ID: %s", userID)
	}

	return nil
}

func (uq *userQuery) GetUserID(email string) (string, error) {
	result := User{}
	query := uq.db.Table("users").Where("email = ?", email).First(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("user record not found, invalid email")
		return "", errors.New("invalid email")
	}

	return result.UserID, nil
}

// GetByID implements user.UserData.
func (uq *userQuery) GetByID(userID string) (user.UserCore, error) {
	var userModel User

	query := uq.db.Preload("Venues").Preload("Reservations").Preload("Reviews").Where("user_id = ?", userID).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			log.Sugar().Errorf("no user found with ID: %s", userID)
			return user.UserCore{}, fmt.Errorf("user not found with ID: %s", userID)
		}
		log.Sugar().Errorf("failed to query user: %w", query.Error)
		return user.UserCore{}, fmt.Errorf("failed to query user: %w", query.Error)
	}

	userCore := UserModelToCore(userModel)
	return userCore, nil
}

// UpdateByID implements user.UserData.
func (uq *userQuery) UpdateByID(userID string, updatedUser user.UserCore) error {
	var userModel User

	// Retrieve the user from the database
	query := uq.db.Where("user_id = ?", userID).First(&userModel)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			log.Sugar().Errorf("no user found with ID: %s", userID)
			return fmt.Errorf("user not found with ID: %s", userID)
		}
		log.Sugar().Errorf("failed to query user: %w", query.Error)
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
			log.Sugar().Errorf("failed to update user: %w", update.Error)
			return fmt.Errorf("failed to update user: %w", update.Error)
		} else if strings.Contains(update.Error.Error(), "'users.email'") {
			log.Error("email already in use")
			return errors.New("email already in use")
		} else if strings.Contains(update.Error.Error(), "'users.phone'") {
			log.Error("phone already in use")
			return errors.New("phone already in use")
		}
	}

	return nil
}
