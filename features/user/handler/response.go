package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/user"
)

type GetUserResponse struct {
	UserID         string
	FullName       string
	Email          string
	Phone          string
	Role           string
	Bio            string
	Address        string
	ProfilePicture string
}

func UserEntityToGetUserResponse(u user.UserEntity) GetUserResponse {
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
