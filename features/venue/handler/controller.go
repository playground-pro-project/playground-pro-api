package handler

import (
	"net/http"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
)

var log = middlewares.Log()

type venueHandler struct {
	service venue.VenueService
}

func New(vs venue.VenueService) venue.VenueHandler {
	return &venueHandler{
		service: vs,
	}
}

// RegisterVenue implements venue.VenueHandler.
func (vh *venueHandler) RegisterVenue() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := RegisterVenueRequest{}
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		errBind := c.Bind(&request)
		if errBind != nil {
			log.Error("error on bind input")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request", nil, nil))
		}

		_, err := vh.service.RegisterVenue(userId, RequestToCore(request))
		if err != nil {
			if strings.Contains(err.Error(), "empty") {
				log.Error("request cannot be empty")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request", nil, nil))
			}
			if strings.Contains(err.Error(), "duplicated") {
				log.Error("error on bind input")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request", nil, nil))
			}
			log.Error("internal server error")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
		}
		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Successfully created new venue", nil, nil))
	}
}

// SearchVenue implements venue.VenueHandler.
func (vh *venueHandler) SearchVenues() echo.HandlerFunc {
	return func(c echo.Context) error {
		var page pagination.Pagination
		limitInt, err := strconv.Atoi(c.QueryParam("limit"))
		if err != nil {
			log.Error("error while converting 'limit' to integer")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad Request", nil, nil))
		}

		pageInt, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil {
			log.Error("error while converting 'page' to integer")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad Request", nil, nil))
		}

		page.Limit = limitInt
		page.Page = pageInt
		page.Sort = c.QueryParam("sort")
		keyword := c.QueryParam("keyword")

		venues, rows, pages, err := vh.service.SearchVenues(keyword, page)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("list venues not found")
				return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "The requested resource was not found", nil, nil))
			} else {
				log.Error("internal server error")
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
			}
		}

		result := make([]SearchVenueResponse, len(venues))
		for i, venue := range venues {
			result[i] = SearchVenue(venue)
		}

		pagination := &pagination.Pagination{
			Limit:      page.Limit,
			Offset:     page.Offset,
			Page:       page.Page,
			TotalRows:  rows,
			TotalPages: pages,
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, pagination))
	}
}

// SelectVenue implements venue.VenueHandler.
func (vh *venueHandler) SelectVenue() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := middlewares.ExtractToken(c)
		if err != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		venueId := c.Param("venue_id")
		venue, err := vh.service.SelectVenue(venueId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("venue not found")
				return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "The requested resource was not found", nil, nil))
			}
			log.Error("internal server error")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
		}

		resp := SelectVenue(venue)
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully operation.", resp, nil))
	}
}

// EditVenue implements venue.VenueHandler.
func (vh *venueHandler) EditVenue() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := EditVenueRequest{}
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		errBind := c.Bind(&request)
		if errBind != nil {
			log.Error("error on bind input")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request", nil, nil))
		}

		venueId := c.Param("venue_id")
		err := vh.service.EditVenue(userId, venueId, RequestToCore(&request))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("venue not found")
				return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "The requested resource was not found", nil, nil))
			}
			log.Error("internal server error")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
		}
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Venue updated successfully", nil, nil))
	}
}

// UnregisterVenue implements venue.VenueHandler.
func (vh *venueHandler) UnregisterVenue() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		venueId := c.Param("venue_id")
		err := vh.service.UnregisterVenue(userId, venueId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("venue not found")
				return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "The requested resource was not found", nil, nil))
			}
			log.Error("internal server error")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
		}
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Venue deleted successfully", nil, nil))
	}
}

// VenueAvailability implements venue.VenueHandler.
func (vh *venueHandler) VenueAvailability() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		venueId := c.Param("venue_id")
		availables, err := vh.service.VenueAvailability(venueId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Info("there is no reservation at this moment")
				return c.JSON(http.StatusFound, helper.ResponseFormat(http.StatusFound, "There is no reservation at this moment", nil, nil))
			}
			log.Error("internal server error")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
		}

		resp := Availability(availables)
		log.Sugar().Infoln(resp)
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully operation.", resp, nil))
	}
}
