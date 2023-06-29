package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/review"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

var log = middlewares.Log()

type reviewHandler struct {
	reviewService review.ReviewService
}

func New(s review.ReviewService) *reviewHandler {
	return &reviewHandler{
		reviewService: s,
	}
}

func (rh *reviewHandler) CreateReview(c echo.Context) error {
	userId, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}

	venueID := c.Param("venue_id")
	req := CreateReviewRequest{}
	errBind := c.Bind(&req)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Invalid request payload"))
	}

	reviewCore := CreateReviewRequestToCore(req)
	_, err := rh.reviewService.CreateReview(venueID, userId, reviewCore)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Review created successfully"))
}

func (rh *reviewHandler) DeleteReview(c echo.Context) error {
	_, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}
	reviewID := c.Param("review_id")
	err := rh.reviewService.DeleteByID(reviewID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "Review deleted successfully"))
}

func (rh *reviewHandler) GetAllReview(c echo.Context) error {
	_, errToken := middlewares.ExtractToken(c)
	if errToken != nil {
		log.Error("missing or malformed JWT")
		return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
	}
	venueID := c.Param("venue_id")
	reviews, err := rh.reviewService.GetAllByVenueID(venueID)
	if err != nil {
		if strings.Contains(err.Error(), "review not found") {
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
	}

	if len(reviews) == 0 {
		log.Error("reservation history not found")
		return helper.NotFoundError(c, "The requested resource was not found")
	}

	reviewsResponse := []GetAllReviewResponse{}
	for _, review := range reviews {
		reviewsResponse = append(reviewsResponse, ReviewCoreToGetAllReviewResponse(review))
	}

	return c.JSON(http.StatusOK, helper.SuccessResponse(reviewsResponse, "Reviews retrieved successfully"))
}
