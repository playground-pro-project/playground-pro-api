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

// SearchVenue implements venue.VenueHandler.
func (vh *venueHandler) SearchVenue() echo.HandlerFunc {
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

		venues, rows, pages, err := vh.service.SearchVenue(keyword, page)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("list venues not found")
				return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "The requested resource was not found", nil, nil))
			} else {
				log.Error("internal server error")
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
			}
		}

		result := make([]searchVenueResponse, len(venues))
		for i, venue := range venues {
			result[i] = searchVenue(venue)
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
