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

func UserRequestToEntity(u RegisterRequest) user.UserEntity {
	return user.UserEntity{
		Fullname: u.FullName,
		Email:    u.Email,
		Phone:    u.Phone,
		Password: u.Password,
	}
}
