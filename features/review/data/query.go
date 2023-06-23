package data

import (
	"errors"
	"fmt"

	"github.com/playground-pro-project/playground-pro-api/features/review"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"gorm.io/gorm"
)

type reviewQuery struct {
	db *gorm.DB
}

// GetAllByVenueID implements review.ReviewData.
func (rq reviewQuery) GetAllByVenueID(venueID string) ([]review.ReviewCore, error) {
	var reviewModels []Review
	query := rq.db.Preload("User").Preload("Venue").Where("venue_id = ?", venueID).Find(&reviewModels)
	if query.Error != nil {
		if errors.Is(query.Error, gorm.ErrRecordNotFound) {
			return []review.ReviewCore{}, fmt.Errorf("review not found with venue ID: %s", venueID)
		}
		return []review.ReviewCore{}, fmt.Errorf("failed to query review: %w", query.Error)
	}

	var reviewCores []review.ReviewCore
	for _, review := range reviewModels {
		reviewCores = append(reviewCores, ReviewModelToCore(review))
	}

	return reviewCores, nil
}

// Create implements review.ReviewData.
func (rq reviewQuery) Create(venueID string, userID string, review review.ReviewCore) (string, error) {
	reviewModel := ReviewCoreToModel(review)
	reviewModel.ReviewID = helper.GenerateReviewID()
	reviewModel.VenueID = venueID
	reviewModel.UserID = userID

	createResult := rq.db.Create(&reviewModel)
	if createResult.Error != nil {
		return "", createResult.Error
	}

	if createResult.RowsAffected == 0 {
		return "", errors.New("failed to insert, row affected is 0")
	}

	return reviewModel.ReviewID, nil
}

// DeleteByID implements review.ReviewData.
func (rq reviewQuery) DeleteByID(reviewID string) error {
	deleteResult := rq.db.Where("review_id = ?", reviewID).Delete(&Review{})
	if deleteResult.Error != nil {
		return fmt.Errorf("failed to delete review: %w", deleteResult.Error)
	}
	if deleteResult.RowsAffected == 0 {
		return fmt.Errorf("no review found with ID: %s", reviewID)
	}

	return nil
}

func New(db *gorm.DB) review.ReviewData {
	return reviewQuery{
		db: db,
	}
}
