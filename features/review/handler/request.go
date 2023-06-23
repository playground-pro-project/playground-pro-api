package handler

import "github.com/playground-pro-project/playground-pro-api/features/review"

type CreateReviewRequest struct {
	// UserID  string  `json:"user_id" form:"user_id"`
	// VenueID string  `json:"venue_id" form:"venue_id"`
	Review string  `json:"review" form:"review"`
	Rating float64 `json:"rating" form:"rating"`
}

func CreateReviewRequestToCore(cr CreateReviewRequest) review.ReviewCore {
	return review.ReviewCore{
		// UserID:  cr.UserID,
		// VenueID: cr.VenueID,
		Review: cr.Review,
		Rating: cr.Rating,
	}
}
