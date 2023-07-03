package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/playground-pro-project/playground-pro-api/features/review"
	"github.com/playground-pro-project/playground-pro-api/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllByVenueID(t *testing.T) {
	data := mocks.NewReviewData(t)
	service := New(data)

	t.Run("success", func(t *testing.T) {
		venueID := "venue_id_1"
		expectedReviewCores := []review.ReviewCore{
			{ReviewID: "review_id_1", Rating: 4},
			{ReviewID: "review_id_2", Rating: 5},
		}
		data.On("GetAllByVenueID", venueID).Return(expectedReviewCores, nil)

		reviewCores, err := service.GetAllByVenueID(venueID)
		assert.NoError(t, err)
		assert.Equal(t, expectedReviewCores, reviewCores)
		data.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		venueID := "venue_id_1"
		expectedErr := errors.New("database error")
		data.On("GetAllByVenueID", venueID).Return(nil, expectedErr)

		reviewCores, err := service.GetAllByVenueID(venueID)
		assert.Error(t, err)
		assert.Equal(t, []review.ReviewCore{}, reviewCores)
		assert.EqualError(t, err, fmt.Sprintf("error: %v", expectedErr))
		data.AssertExpectations(t)
	})
}

func TestCreateReview(t *testing.T) {
	data := &mocks.ReviewData{}
	service := New(data)

	t.Run("success", func(t *testing.T) {
		venueID := "venue_id_1"
		userID := "user_id_1"
		reviewCore := review.ReviewCore{Rating: 4}
		expectedReviewID := "review_id_1"
		data.On("Create", venueID, userID, reviewCore).Return(expectedReviewID, nil)

		reviewID, err := service.CreateReview(venueID, userID, reviewCore)
		assert.NoError(t, err)
		assert.Equal(t, expectedReviewID, reviewID)
		data.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		venueID := "venue_id_1"
		userID := "user_id_1"
		reviewCore := review.ReviewCore{Rating: 4}
		expectedErr := errors.New("database error")
		data.On("Create", venueID, userID, reviewCore).Return("", expectedErr)

		reviewID, err := service.CreateReview(venueID, userID, reviewCore)
		assert.Error(t, err)
		assert.Equal(t, "", reviewID)
		assert.EqualError(t, err, fmt.Sprintf("%v", expectedErr))
		data.AssertExpectations(t)
	})

	t.Run("nil reviewID", func(t *testing.T) {
		venueID := "venue_id_1"
		userID := "user_id_1"
		reviewCore := review.ReviewCore{Rating: 4}
		data.On("Create", venueID, userID, reviewCore).Return(nil, nil)

		reviewID, err := service.CreateReview(venueID, userID, reviewCore)
		assert.Error(t, err)
		assert.Equal(t, "", reviewID)
		assert.EqualError(t, err, "failed to create review: reviewID is nil")
		data.AssertExpectations(t)
	})
}

func TestDeleteReview(t *testing.T) {
	data := &mocks.ReviewData{}
	service := New(data)

	t.Run("success", func(t *testing.T) {
		reviewID := "review_id_1"
		data.On("DeleteByID", reviewID).Return(nil)

		err := service.DeleteByID(reviewID)
		assert.NoError(t, err)
		data.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		reviewID := "review_id_1"
		expectedErr := errors.New("database error")
		data.On("DeleteByID", reviewID).Return(expectedErr)

		err := service.DeleteByID(reviewID)
		assert.Error(t, err)
		assert.EqualError(t, err, fmt.Sprintf("error: %v", expectedErr))
		data.AssertExpectations(t)
	})
}
