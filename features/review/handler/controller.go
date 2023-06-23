package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/review"
)

type reviewHandler struct {
	reviewService review.ReviewService
}

func New(s review.ReviewService) *reviewHandler {
	return &reviewHandler{
		reviewService: s,
	}
}

func (rh *reviewHandler) CreateReview(c echo.Context) error {
	userID := middlewares.ExtractUserIDFromToken(c)
	venueID := c.Param("venue_id")

	req := CreateReviewRequest{}
	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}

	reviewCore := CreateReviewRequestToCore(req)
	_, err = rh.reviewService.CreateReview(venueID, userID, reviewCore)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Review created successfully",
	})
}

func (rh *reviewHandler) DeleteReview(c echo.Context) error {
	// userID := middlewares.ExtractUserIDFromToken(c)
	reviewID := c.Param("review_id")
	err := rh.reviewService.DeleteByID(reviewID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": "Review deleted successfully",
	})
}
