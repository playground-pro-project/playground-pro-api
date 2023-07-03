package service

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/mocks"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// func TestRegister(t *testing.T) {
// 	data := mocks.NewUserData(t)
// 	validate := validator.New()
// 	arguments := user.UserCore{
// 		Fullname: "admin",
// 		Email:    "admin@gmail.com",
// 		Phone:    "081235288543",
// 		Password: "@S3#cr3tP4ss#word123",
// 		Role:     "user",
// 	}
// 	result := user.UserCore{
// 		UserID:   "550e8400-e29b-41d4-a716-446655440000",
// 		Fullname: "admin",
// 		Email:    "admin@gmail.com",
// 		Phone:    "081235288543",
// 		Password: "@S3#cr3tP4ss#word123",
// 		Role:     "user",
// 	}
// 	service := New(data, validate)

// 	t.Run("fullname cannot be empty", func(t *testing.T) {
// 		request := user.UserCore{
// 			Fullname: "",
// 			Email:    "admin@gmail.com",
// 			Phone:    "081235288543",
// 			Password: "@S3#cr3tP4ss#word123",
// 			Role:     "user",
// 		}
// 		_, _, err := service.Register(request)
// 		expectedErr := errors.New("fullname is required")
// 		assert.NotNil(t, err)
// 		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("invalid email format", func(t *testing.T) {
// 		request := user.UserCore{
// 			Fullname: "admin",
// 			Email:    "admin@.com",
// 			Phone:    "081235288543",
// 			Password: "@S3#cr3tP4ss#word123",
// 			Role:     "user",
// 		}
// 		_, _, err := service.Register(request)
// 		expectedErr := errors.New("wrong email format")
// 		assert.NotNil(t, err)
// 		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("phone cannot be empty", func(t *testing.T) {
// 		request := user.UserCore{
// 			Fullname: "admin",
// 			Email:    "admin@gmail.com",
// 			Phone:    "",
// 			Password: "@S3#cr3tP4ss#word123",
// 			Role:     "user",
// 		}
// 		_, _, err := service.Register(request)
// 		expectedErr := errors.New("phone is required")
// 		assert.NotNil(t, err)
// 		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("password cannot be empty", func(t *testing.T) {
// 		request := user.UserCore{
// 			Fullname: "admin",
// 			Email:    "admin@gmail.com",
// 			Phone:    "081235288543",
// 			Password: "",
// 			Role:     "user",
// 		}
// 		_, _, err := service.Register(request)
// 		expectedErr := errors.New("password should be at least 6 characters long")
// 		assert.NotNil(t, err)
// 		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("success create account", func(t *testing.T) {
// 		data.On("Register", mock.Anything).Return(result, nil).Once()
// 		res, _, err := service.Register(arguments)
// 		assert.Nil(t, err)
// 		assert.Equal(t, result.UserID, res.UserID)
// 		assert.NotEmpty(t, result.Fullname)
// 		assert.NotEmpty(t, result.Email)
// 		assert.NotEmpty(t, result.Password)
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("error while hashing password", func(t *testing.T) {
// 		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("error while hashing password")).Once()
// 		res, _, err := service.Register(arguments)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "", res.UserID)
// 		assert.ErrorContains(t, err, "password")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("error insert data, duplicated", func(t *testing.T) {
// 		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("error insert data, duplicated")).Once()
// 		res, _, err := service.Register(arguments)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "", res.UserID)
// 		assert.ErrorContains(t, err, "duplicated")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("internal server error", func(t *testing.T) {
// 		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("server error")).Once()
// 		res, _, err := service.Register(arguments)
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "", res.UserID)
// 		assert.ErrorContains(t, err, "server error")
// 		data.AssertExpectations(t)
// 	})

// 	t.Run("user already exists", func(t *testing.T) {
// 		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("user already exists")).Once()
// 		_, _, err := service.Register(arguments)
// 		assert.NotNil(t, err)
// 		assert.ErrorContains(t, err, "user already exists")
// 		data.AssertExpectations(t)
// 	})
// }

func TestLogin(t *testing.T) {
	data := mocks.NewUserData(t)
	arguments := user.UserCore{Email: "admin@gmail.com", Password: "@SecretPassword123"}
	wrongArguments := user.UserCore{Email: "admin@gmail.com", Password: "@WrongPassword"}
	token := "123"
	emptyToken := ""
	hashed, _ := helper.HashPassword(arguments.Password)
	result := user.UserCore{UserID: "uuid", Fullname: "admin", Password: hashed}
	validate := validator.New()
	service := New(data, validate)

	t.Run("invalid email format", func(t *testing.T) {
		request := user.UserCore{
			Email:    "admin@.com",
			Password: "@S3#cr3tP4ss#word123",
		}
		_, _, err := service.Login(request)
		expectedErr := errors.New("invalid email format")
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
		data.AssertExpectations(t)
	})

	t.Run("password cannot be empty", func(t *testing.T) {
		request := user.UserCore{
			Email:    "admin@gmail.com",
			Password: "",
		}
		_, _, err := service.Login(request)
		expectedErr := errors.New("password cannot be empty")
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
		data.AssertExpectations(t)
	})

	t.Run("success login", func(t *testing.T) {
		data.On("Login", mock.Anything).Return(result, token, nil).Once()
		res, token, err := service.Login(arguments)
		assert.Nil(t, err)
		assert.NotEmpty(t, token)
		assert.Equal(t, result.Email, res.Email)
		assert.Equal(t, result.Password, res.Password)
		data.AssertExpectations(t)
	})

	t.Run("invalid email and password", func(t *testing.T) {
		data.On("Login", mock.Anything).Return(user.UserCore{}, token, errors.New("invalid email and password")).Once()
		_, _, err := service.Login(wrongArguments)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "invalid email and password")
		data.AssertExpectations(t)
	})

	t.Run("password does not match", func(t *testing.T) {
		data.On("Login", mock.Anything).Return(user.UserCore{}, token, errors.New("password does not match")).Once()
		_, _, err := service.Login(wrongArguments)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "password does not match")
		data.AssertExpectations(t)
	})

	t.Run("error while creating jwt token", func(t *testing.T) {
		data.On("Login", mock.Anything).Return(user.UserCore{}, token, errors.New("error while creating jwt token")).Once()
		_, _, err := service.Login(wrongArguments)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error while creating jwt token")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Login", mock.Anything).Return(user.UserCore{}, emptyToken, errors.New("server error")).Once()
		res, token, err := service.Login(arguments)
		assert.NotNil(t, err)
		assert.Equal(t, "", res.UserID)
		assert.Equal(t, emptyToken, token)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestDeleteByID(t *testing.T) {
	data := new(mocks.UserData)
	validator := new(validator.Validate)
	service := New(data, validator)
	userID := "user_id_1"
	t.Run("success", func(t *testing.T) {
		data.On("DeleteByID", userID).Return(nil).Once()
		err := service.DeleteByID(userID)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockError := errors.New("error deleting user")
		data.On("DeleteByID", userID).Return(mockError).Once()
		err := service.DeleteByID(userID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error deleting user")
		data.AssertExpectations(t)
	})
}

func TestGetByID(t *testing.T) {
	data := new(mocks.UserData)
	validator := new(validator.Validate)
	service := New(data, validator)

	userID := "user_id_1"

	t.Run("success", func(t *testing.T) {
		mockUser := user.UserCore{
			UserID:   userID,
			Fullname: "John Doe",
			Email:    "johndoe@example.com",
			Phone:    "123456789",
		}

		data.On("GetByID", userID).Return(mockUser, nil).Once()
		result, err := service.GetByID(userID)
		assert.Nil(t, err)
		assert.Equal(t, mockUser, result)
		data.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockError := errors.New("error getting user")
		data.On("GetByID", userID).Return(user.UserCore{}, mockError).Once()
		_, err := service.GetByID(userID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error getting user")
		data.AssertExpectations(t)
	})
}

func TestGetUserID(t *testing.T) {
	data := new(mocks.UserData)
	validator := new(validator.Validate)
	service := New(data, validator)

	email := "johndoe@example.com"
	userID := "user_id_1"

	t.Run("success", func(t *testing.T) {
		data.On("GetUserID", email).Return(userID, nil).Once()
		result, err := service.GetUserID(email)
		assert.Nil(t, err)
		assert.Equal(t, userID, result)
		data.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		mockError := errors.New("error getting user ID")
		data.On("GetUserID", email).Return("", mockError).Once()
		_, err := service.GetUserID(email)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error getting user ID")
		data.AssertExpectations(t)
	})
}

func TestUpdateByID(t *testing.T) {
	data := new(mocks.UserData)
	validator := validator.New()
	service := New(data, validator)

	userID := "user_id_1"
	updatedUser := user.UserCore{
		UserID:   userID,
		Fullname: "Updated Name",
		Email:    "updated@example.com",
		Phone:    "987654321",
		Password: "newpassword",
	}

	t.Run("valid update", func(t *testing.T) {
		data.On("UpdateByID", userID, updatedUser).Return(nil).Once()
		err := service.UpdateByID(userID, updatedUser)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("invalid email format", func(t *testing.T) {
		updatedUser.Email = "invalid_email"
		err := service.UpdateByID(userID, updatedUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "invalid email format")
		data.AssertExpectations(t)
	})

	t.Run("invalid password format", func(t *testing.T) {
		updatedUser.Password = "short" // Use an invalid password, e.g., too short.
		err := service.UpdateByID(userID, updatedUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "password should contain at least one special character")
		data.AssertExpectations(t)
	})

	t.Run("error updating user", func(t *testing.T) {
		mockError := errors.New("error updating user")
		data.On("UpdateByID", userID, updatedUser).Return(mockError).Once()
		err := service.UpdateByID(userID, updatedUser)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error updating user")
		data.AssertExpectations(t)
	})
}
