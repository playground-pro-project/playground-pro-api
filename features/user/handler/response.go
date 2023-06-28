package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/user"
)

type RegisterResponse struct {
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
	OTP    string `json:"otp,omitempty"`
}

type LoginResponse struct {
	UserID        string `json:"user_id,omitempty"`
	Email         string `json:"email,omitempty"`
	Token         string `json:"token,omitempty"`
	Role          string `json:"role,omitempty"`
	AccountStatus string `json:"account_status,omitempty"`
}

type GetUserResponse struct {
	UserID         string `json:"user_id,omitempty"`
	FullName       string `json:"fullname,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Role           string `json:"role,omitempty"`
	Bio            string `json:"bio,omitempty"`
	Address        string `json:"address,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
}

func UserCoreToLoginResponse(u user.UserCore) LoginResponse {
	return LoginResponse{
		UserID:        u.UserID,
		Email:         u.Email,
		Role:          u.Role,
		AccountStatus: u.AccountStatus,
	}
}

func UserCoreToGetUserResponse(u user.UserCore) GetUserResponse {
	return GetUserResponse{
		UserID:         u.UserID,
		FullName:       u.Fullname,
		Email:          u.Email,
		Phone:          u.Phone,
		Role:           u.Role,
		Bio:            u.Bio,
		Address:        u.Address,
		ProfilePicture: u.ProfilePicture,
	}
}

func UserCoreToRegisterResponse(u user.UserCore) RegisterResponse {
	return RegisterResponse{
		UserID: u.UserID,
		Email:  u.Email,
	}
}
