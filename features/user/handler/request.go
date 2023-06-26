package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/user"
)

const (
	maxFileSize              = 1 << 20 // 1 MB
	profilePictureBaseURL    = "https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/user-profile-picture/"
	defaultProfilePictureURL = "https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png"
	ownerFileBaseURL         = "https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/owner-docs/"
)

type RegisterRequest struct {
	FullName string `json:"full_name" form:"full_name"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Password string `json:"password" form:"password"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type OTPInput struct {
	UserID  string `json:"user_id"`
	OTPCode string `json:"otp_code"`
}

type EditProfileRequest struct {
	Fullname string `json:"full_name" form:"full_name"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Address  string `json:"address" form:"address"`
	Bio      string `json:"bio" form:"bio"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" form:"old_password"`
	NewPassword string `json:"new_password" form:"new_password"`
}

func RequestToCore(data interface{}) user.UserCore {
	res := user.UserCore{}
	switch v := data.(type) {
	case RegisterRequest:
		res.Fullname = v.FullName
		res.Email = v.Email
		res.Phone = v.Phone
		res.Password = v.Password
	case LoginRequest:
		res.Email = v.Email
		res.Password = v.Password
	case OTPInput:
		res.UserID = v.UserID
		res.OTPCode = v.OTPCode
	default:
		return user.UserCore{}
	}
	return res
}

func EditProfileRequestToCore(e EditProfileRequest) user.UserCore {
	return user.UserCore{
		Fullname: e.Fullname,
		Email:    e.Email,
		Phone:    e.Phone,
		Address:  e.Address,
		Bio:      e.Bio,
	}
}

func UpdatePasswordRequestToCore(p ChangePasswordRequest) user.UserCore {
	return user.UserCore{
		Password: p.NewPassword,
	}
}
