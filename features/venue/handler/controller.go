package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/aws"
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
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		errBind := c.Bind(&request)
		if errBind != nil {
			log.Error("error on bind input")
			return helper.BadRequestError(c, "Bad request")
		}

		_, err := vh.service.RegisterVenue(userId, RequestToCore(request))
		if err != nil {
			if strings.Contains(err.Error(), "empty") {
				log.Error("request cannot be empty")
				return helper.BadRequestError(c, "Bad request")
			}
			if strings.Contains(err.Error(), "error insert data, duplicated") {
				log.Error("error insert data, duplicated")
				return helper.BadRequestError(c, "Bad request")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
		}
		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Successfully operation", nil, nil))
	}
}

// SearchVenue implements venue.VenueHandler.
func (vh *venueHandler) SearchVenues() echo.HandlerFunc {
	return func(c echo.Context) error {
		var page pagination.Pagination
		limitInt, _ := strconv.Atoi(c.QueryParam("limit"))
		pageInt, _ := strconv.Atoi(c.QueryParam("page"))
		page.Limit = limitInt
		page.Page = pageInt
		page.Sort = c.QueryParam("sort")
		keyword := c.QueryParam("keyword")

		venues, rows, pages, err := vh.service.SearchVenues(keyword, page)
		if err != nil {
			if strings.Contains(err.Error(), "venues not found") {
				log.Error("venues not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
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

		if len(result) == 0 {
			log.Error("venues not found")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, pagination))
	}
}

// SelectVenue implements venue.VenueHandler.
func (vh *venueHandler) SelectVenue() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venueId := c.Param("venue_id")
		if venueId == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		venue, err := vh.service.SelectVenue(venueId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("venue not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
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
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		errBind := c.Bind(&request)
		if errBind != nil {
			log.Error("error on bind input")
			return helper.BadRequestError(c, "Bad request")
		}

		venueId := c.Param("venue_id")
		if venueId == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		err := vh.service.EditVenue(userId, venueId, RequestToCore(&request))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("venue not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else if strings.Contains(err.Error(), "no venue has been created") {
				log.Error("no venue has been created")
				return helper.NotFoundError(c, "no venue has been created")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
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
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venueId := c.Param("venue_id")
		if venueId == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		err := vh.service.UnregisterVenue(userId, venueId)
		if err != nil {
			if strings.Contains(err.Error(), "venue record not found") {
				log.Error("venue record not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else if strings.Contains(err.Error(), "no row affected") {
				log.Error("no row affected")
				return helper.NotFoundError(c, "The requested resource was not found")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
		}
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully deleted a venue", nil, nil))
	}
}

// VenueAvailability implements venue.VenueHandler.
func (vh *venueHandler) VenueAvailability() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venueId := c.Param("venue_id")
		if venueId == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		availables, err := vh.service.VenueAvailability(venueId)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Info("there is no reservation at this moment")
				return helper.NotFoundError(c, "The requested resource was not found")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
		}

		resp := Availability(availables)
		log.Sugar().Infoln(resp)
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully operation.", resp, nil))
	}
}

func (vh *venueHandler) CreateVenueImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venueId := c.Param("venue_id")
		if venueId == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		form, err := c.MultipartForm()
		if err != nil {
			log.Error("Failed to retrieve file: " + err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Failed to retrieve file: "+err.Error()))
		}

		files := form.File["files"]
		for _, file := range files {
			// Get the file type from the Content-Type header
			fileType := file.Header.Get("Content-Type")
			fileExtension := filepath.Ext(file.Filename)
			fileExtension = strings.ToLower(fileExtension)

			allowedExtensions := []string{".jpg", ".jpeg", ".png"}
			extensionAllowed := false
			for _, ext := range allowedExtensions {
				if ext == fileExtension {
					extensionAllowed = true
					break
				}
			}
			if !extensionAllowed {
				log.Error(fileExtension + " is not allowed")
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse(fileExtension+" is not allowed. Only JPG, JPEG, and PNG files are allowed."))
			}

			fileSize := file.Size
			if fileSize > maxVenueFileSize {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("Please upload a file smaller than 2 MB."))
			}

			id := helper.GenerateIdentifier()
			filename := id + "-" + file.Filename
			path := "venue-images/" + filename

			fileContent, err := file.Open()
			if err != nil {
				log.Error("Failed to open file: " + err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to open file: "+err.Error()))
			}
			defer fileContent.Close()

			image := venue.VenuePictureCore{
				VenueID: venueId,
				URL:     fmt.Sprintf("%s%s", venueFileBaseURL, filepath.Base(filename)),
			}

			_, err = vh.service.CreateVenueImage(image)
			if err != nil {
				log.Error("Failed to insert image. " + err.Error())
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("Failed to insert image. "+err.Error()))
			}

			awsService := aws.InitS3()

			// Upload profile picture file to cloud
			err = awsService.UploadFile(path, fileType, fileContent)
			if err != nil {
				log.Error("Failed to upload file to cloud service: " + err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to upload file to cloud service: "+err.Error()))
			}
		}

		log.Sugar().Infof(venueId + " venue image added successfully")
		return c.JSON(http.StatusCreated, helper.SuccessResponse(nil, "venue image added successfully"))
	}
}

func (vh *venueHandler) DeleteVenueImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venueID := c.Param("venue_id")
		if venueID == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		venueImageId := c.Param("image_id")
		if venueImageId == "" {
			log.Error("empty image_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		vn, err := vh.service.GetVenueImageByID(venueID, venueImageId)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusNotFound, helper.ErrorResponse(err.Error()))
		}

		// Delete profile picture in the cloud before updated
		awsService := aws.InitS3()

		prevFilename := filepath.Base(vn.URL)
		prevPath := "venue-images/" + prevFilename

		err = awsService.DeleteFile(prevPath)
		if err != nil {
			log.Error("Failed to delete file from cloud service: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to delete file from cloud service: "+err.Error()))
		}

		err = vh.service.DeleteVenueImage(venueImageId)
		if err != nil {
			if strings.Contains(err.Error(), "failed to delete") {
				log.Error("failed to delete image")
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to delete image"))
			} else {
				log.Error("Internal server error")
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
			}
		}

		log.Sugar().Infof(venueImageId + " venue image deleted successfully")
		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "venue image deleted successfully"))
	}
}

func (vh *venueHandler) GetAllVenueImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		venueId := c.Param("venue_id")
		if venueId == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		images, err := vh.service.GetAllVenueImage(venueId)
		if err != nil {
			if strings.Contains(err.Error(), "no images found") {
				log.Error("no images found for venue")
				return c.JSON(http.StatusNotFound, helper.ErrorResponse("no images found for venue"))
			}
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
		}

		var resp []GetAllVenueImageResponse
		for _, image := range images {
			resp = append(resp, GetAllVenueImageToResponse(image))
		}

		log.Sugar().Infoln(resp)
		return c.JSON(http.StatusOK, helper.SuccessResponse(resp, "venue image retrieved successfully"))
	}
}

// MyVenues implements venue.VenueHandler.
func (vh *venueHandler) MyVenues() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venues, err := vh.service.MyVenues(userId)
		if err != nil {
			if strings.Contains(err.Error(), "venues not found") {
				log.Error("venues not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		result := make([]SearchVenueResponse, len(venues))
		for i, venue := range venues {
			result[i] = SearchVenue(venue)
		}

		if len(result) == 0 {
			log.Error("venues not found")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, nil))
	}
}
