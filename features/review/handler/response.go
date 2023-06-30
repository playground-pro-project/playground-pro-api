package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/review"
)

type GetAllReviewResponse struct {
	ReviewID string       `json:"review_id`
	UserID   string       `json:"user_id"`
	Review   string       `json:"review"`
	Rating   float64      `json:"rating"`
	User     UserResponse `json:"user"`
}

type UserResponse struct {
	Fullname string `json:"fullname"`
}

func ReviewCoreToGetAllReviewResponse(r review.ReviewCore) GetAllReviewResponse {
	return GetAllReviewResponse{
		ReviewID: r.ReviewID,
		UserID:   r.UserID,
		Review:   r.Review,
		Rating:   r.Rating,
		User:     UserCoreToUserResponse(r.User),
	}
}

func UserCoreToUserResponse(u review.UserCore) UserResponse {
	return UserResponse{
		Fullname: u.Fullname,
	}
}
