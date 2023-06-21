package data

import (
	"errors"
	"strings"

	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"gorm.io/gorm"
)

type userQuery struct {
	db *gorm.DB
}

// Create implements user.UserData.
func (uq userQuery) Create(user user.UserEntity) (string, error) {
	userModel := UserEntityToModel(user)
	userModel.Password = helper.HashPass(userModel.Password)
	userModel.UserID = helper.GenerateUserID()

	createResult := uq.db.Create(&userModel)
	if createResult.Error != nil {
		switch {
		case strings.Contains(createResult.Error.Error(), "email"):
			return "", errors.New("email already in use")
		case strings.Contains(createResult.Error.Error(), "phone"):
			return "", errors.New("phone already in use")
		}
		return "", createResult.Error
	}

	if createResult.RowsAffected == 0 {
		return "", errors.New("failed to insert, row affected is 0")
	}

	return userModel.UserID, nil
}

// DeleteByID implements user.UserData.
func (uq userQuery) DeleteByID(userID string) error {
	panic("unimplemented")
}

// GetAll implements user.UserData.
func (uq userQuery) GetAll() ([]user.UserEntity, error) {
	panic("unimplemented")
}

// GetByID implements user.UserData.
func (uq userQuery) GetByID(userID string) (user.UserEntity, error) {
	panic("unimplemented")
}

// Login implements user.UserData.
func (uq userQuery) Login(email string, password string) (user.UserEntity, string, error) {
	panic("unimplemented")
}

// UpdateByID implements user.UserData.
func (uq userQuery) UpdateByID(userID string, updatedUser user.UserEntity) error {
	panic("unimplemented")
}

func New(db *gorm.DB) user.UserData {
	return userQuery{
		db: db,
	}
}
