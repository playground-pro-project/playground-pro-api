package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/user"
)

type loginResponse struct {
	UserID     string `json:"user_id"`
	Email      string `json:"email"`
	Token      string `json:"token"`
	OTPEnabled bool   `json:"otp_enabled"`
}

type otpResponse struct {
	UserID      string `json:"user_id,omitempty"`
	OTPSecret   string `json:"otp_secret,omitempty"`
	OTPAuthURL  string `json:"otp_auth_url,omitempty"`
	OTPEnabled  bool   `json:"otp_enabled,omitempty"`
	OTPVerified bool   `json:"otp_verified,omitempty"`
}

type GetUserResponse struct {
	UserID         string `json:"user_id,omitempty"`
	FullName       string `json:"full_name,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Role           string `json:"role,omitempty"`
	Bio            string `json:"bio,omitempty"`
	Address        string `json:"address,omitempty"`
	ProfilePicture string `json:"profile_picture,omitempty"`
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
