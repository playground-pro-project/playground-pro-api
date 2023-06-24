package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/review"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
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
	userId, err := middlewares.ExtractToken(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}
	venueID := c.Param("venue_id")

	req := CreateReviewRequest{}
	errBind := c.Bind(&req)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request payload",
		})
	}

	reviewCore := CreateReviewRequestToCore(req)
	_, err = rh.reviewService.CreateReview(venueID, userId, reviewCore)
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

func (rh *reviewHandler) GetAllReview(c echo.Context) error {
	venueID := c.Param("venue_id")
	reviews, err := rh.reviewService.GetAllByVenueID(venueID)
	if err != nil {
		if strings.Contains(err.Error(), "review not found") {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	reviewsResponse := []GetAllReviewResponse{}
	for _, review := range reviews {
		reviewsResponse = append(reviewsResponse, ReviewCoreToGetAllReviewResponse(review))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data":    reviewsResponse,
		"message": "Reviews retrieved successfully",
	})
}
