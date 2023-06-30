package service

import (
	"fmt"

	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/review"
)

var log = middlewares.Log()

type reviewService struct {
	reviewData review.ReviewData
}

// GetAllByVenueID implements review.ReviewService.
func (rs *reviewService) GetAllByVenueID(venueID string) ([]review.ReviewCore, error) {
	reviewCores, err := rs.reviewData.GetAllByVenueID(venueID)
	if err != nil {
		return []review.ReviewCore{}, fmt.Errorf("error: %w", err)
	}

	return reviewCores, nil
}

// CreateReview implements review.ReviewService.
func (rs *reviewService) CreateReview(venueID string, userID string, review review.ReviewCore) (string, error) {
	reviewID, err := rs.reviewData.Create(venueID, userID, review)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return reviewID, nil
}

// DeleteReview implements review.ReviewService.
func (rs *reviewService) DeleteByID(reviewID string) error {
	err := rs.reviewData.DeleteByID(reviewID)
	if err != nil {
		log.Sugar().Errorf("error: %w", err)
		return fmt.Errorf("error: %w", err)
	}

	log.Sugar().Info("success delete review")
	return nil
}

func New(repo review.ReviewData) review.ReviewService {
	return &reviewService{
		reviewData: repo,
	}
}
