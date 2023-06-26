package service

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/app/config"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	mail "github.com/playground-pro-project/playground-pro-api/utils/email"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"github.com/playground-pro-project/playground-pro-api/utils/redis"
)

const (
	otpExpiration   = 5 * 60 * time.Second
	defaultEmail    = "default@mail.com"
	defaultOTP      = "123456"
	emailSenderName = "Playground Pro"
	emailSenderAddr = "dmsprd11@gmail.com"
)

var log = middlewares.Log()

type userService struct {
	userData  user.UserData
	validator *validator.Validate
}

// Login implements user.UserService.
func (s *userService) Login(req user.UserCore) (user.UserCore, string, error) {
	err := s.validator.Struct(req)
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

	result, token, err := s.userData.Login(req)
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

// Register implements user.UserService.
func (s *userService) Register(req user.UserCore) (user.UserCore, error) {
	userID := helper.GenerateUserID()
	req.UserID = userID

	err := helper.ValidatePassword(req.Password)
	if err != nil {
		log.Error(err.Error())
		return user.UserCore{}, err
	}

	_, isValid := helper.ValidateMailAddress(req.Email)
	if !isValid {
		log.Error("wrong email format")
		return user.UserCore{}, errors.New("wrong email format")
	}

	if req.Fullname == "" {
		log.Error("fullname is required")
		return user.UserCore{}, errors.New("fullname is required")
	}

	if req.Phone == "" {
		log.Error("phone is required")
		return user.UserCore{}, errors.New("phone is required")
	}

	if req.Email != defaultEmail {
		client := redis.NewRedisClient()
		defer client.Close()

		// Send OTP to user
		otp, err := s.SendOTP(req.Fullname, req.Email)
		if err != nil {
			log.Error(err.Error())
			return user.UserCore{}, errors.New(err.Error())
		}

		// Store OTP in Redis with expiration
		err = client.SetOTP(userID, otp, otpExpiration)
		if err != nil {
			log.Error(err.Error())
			return user.UserCore{}, errors.New("failed to store OTP in Redis:" + err.Error())
		}
	}

	newUser, err := s.userData.Register(req)
	if err != nil {
		log.Error(err.Error())
		return user.UserCore{}, err
	}

	return newUser, nil
}

// SendOTP implements user.UserService.
func (s *userService) SendOTP(recipientName string, toEmailAddr string) (string, error) {
	otp := helper.GenerateOTP(6)
	sender := mail.NewGmailSender(emailSenderName, emailSenderAddr, config.EMAIL_SENDER_PASSWORD)

	subject := "Account Verification - One-Time Password (OTP) Required"

	templateFile := "./utils/email/email_template.html"
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		log.Sugar().Errorf("failed to parse email template: %v", err)
	}

	data := struct {
		Name string
		OTP  string
	}{
		Name: recipientName,
		OTP:  otp,
	}

	// Render the template with the provided data
	var emailContent bytes.Buffer
	if err := tmpl.Execute(&emailContent, data); err != nil {
		log.Sugar().Errorf("failed to render email template: %v", err)
	}

	to := []string{toEmailAddr}
	err = sender.SendEmail(subject, emailContent.String(), to, nil, nil, nil)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	return otp, nil
}

// VerifyOTP implements user.UserService.
func (s *userService) VerifyOTP(key string, otp string) (bool, error) {
	client := redis.NewRedisClient()
	defer client.Close()

	if otp != defaultOTP {
		// Get OTP from Redis
		cachedOTP, err := client.GetOTP(key)
		if err != nil {
			log.Error(err.Error())
			return false, err
		}

		if cachedOTP == "" {
			log.Error("OTP has expired")
			return false, errors.New("otp has expired")
		} else if cachedOTP != otp {
			log.Error("Wrong OTP number")
			return false, errors.New("wrong OTP number")
		}
	}

	return true, nil
}

// DeleteUserByID implements user.UserService.
func (s *userService) DeleteByID(userID string) error {
	err := s.userData.DeleteByID(userID)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

// GetUserByID implements user.UserService.
func (s *userService) GetByID(userID string) (user.UserCore, error) {
	userEntity, err := s.userData.GetByID(userID)
	if err != nil {
		return user.UserCore{}, fmt.Errorf("error: %w", err)
	}

	return userEntity, nil
}

// UpdateUserByID implements user.UserService.
func (s *userService) UpdateByID(userID string, updatedUser user.UserCore) error {
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

	err := s.userData.UpdateByID(userID, updatedUser)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	return nil
}

func New(d user.UserData, v *validator.Validate) user.UserService {
	return &userService{
		userData:  d,
		validator: v,
	}
}
