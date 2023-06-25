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

func TestRegister(t *testing.T) {
	data := mocks.NewUserData(t)
	validate := validator.New()
	arguments := user.UserCore{
		Fullname: "admin",
		Email:    "admin@gmail.com",
		Phone:    "081235288543",
		Password: "@S3#cr3tP4ss#word123",
		Role:     "user",
	}
	result := user.UserCore{
		UserID:   "550e8400-e29b-41d4-a716-446655440000",
		Fullname: "admin",
		Email:    "admin@gmail.com",
		Phone:    "081235288543",
		Password: "@S3#cr3tP4ss#word123",
		Role:     "user",
	}
	service := New(data, validate)

	t.Run("fullname cannot be empty", func(t *testing.T) {
		request := user.UserCore{
			Fullname: "",
			Email:    "admin@gmail.com",
			Phone:    "081235288543",
			Password: "@S3#cr3tP4ss#word123",
			Role:     "user",
		}
		_, err := service.Register(request)
		expectedErr := errors.New("fullname cannot be empty")
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
		data.AssertExpectations(t)
	})

	t.Run("invalid email format", func(t *testing.T) {
		request := user.UserCore{
			Fullname: "admin",
			Email:    "admin@.com",
			Phone:    "081235288543",
			Password: "@S3#cr3tP4ss#word123",
			Role:     "user",
		}
		_, err := service.Register(request)
		expectedErr := errors.New("invalid email format")
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
		data.AssertExpectations(t)
	})

	t.Run("phone cannot be empty", func(t *testing.T) {
		request := user.UserCore{
			Fullname: "admin",
			Email:    "admin@gmail.com",
			Phone:    "",
			Password: "@S3#cr3tP4ss#word123",
			Role:     "user",
		}
		_, err := service.Register(request)
		expectedErr := errors.New("phone cannot be empty")
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
		data.AssertExpectations(t)
	})

	t.Run("password cannot be empty", func(t *testing.T) {
		request := user.UserCore{
			Fullname: "admin",
			Email:    "admin@gmail.com",
			Phone:    "081235288543",
			Password: "",
			Role:     "user",
		}
		_, err := service.Register(request)
		expectedErr := errors.New("password cannot be empty")
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErr.Error(), "Expected error message does not match")
		data.AssertExpectations(t)
	})

	t.Run("success create account", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(result, nil).Once()
		res, err := service.Register(arguments)
		assert.Nil(t, err)
		assert.Equal(t, result.UserID, res.UserID)
		assert.NotEmpty(t, result.Fullname)
		assert.NotEmpty(t, result.Email)
		assert.NotEmpty(t, result.Password)
		data.AssertExpectations(t)
	})

	t.Run("error while hashing password", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("error while hashing password")).Once()
		res, err := service.Register(arguments)
		assert.NotNil(t, err)
		assert.Equal(t, "", res.UserID)
		assert.ErrorContains(t, err, "password")
		data.AssertExpectations(t)
	})

	t.Run("error insert data, duplicated", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("error insert data, duplicated")).Once()
		res, err := service.Register(arguments)
		assert.NotNil(t, err)
		assert.Equal(t, "", res.UserID)
		assert.ErrorContains(t, err, "duplicated")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("Register", mock.Anything).Return(user.UserCore{}, errors.New("server error")).Once()
		res, err := service.Register(arguments)
		assert.NotNil(t, err)
		assert.Equal(t, "", res.UserID)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

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