package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/user"
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

func RegisterRequestToCore(r RegisterRequest) user.UserCore {
	return user.UserCore{
		Fullname: r.FullName,
		Email:    r.Email,
		Phone:    r.Phone,
		Password: r.Password,
	}
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
