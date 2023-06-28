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
			if strings.Contains(err.Error(), "error insert data, duplicated") {
				log.Error("error insert data, duplicated")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request", nil, nil))
			}
			log.Error("internal server error")
			return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
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

		if len(result) == 0 {
			log.Error("venues not found")
			return c.JSON(http.StatusNotFound, helper.ResponseFormat(http.StatusNotFound, "The requested resource was not found", nil, nil))
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

func (vh *venueHandler) CreateVenueImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := middlewares.ExtractToken(c)
		if err != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		venueID := c.Param("venue_id")
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

			awsService := aws.InitS3()

			// Upload profile picture file to cloud
			err = awsService.UploadFile(path, fileType, fileContent)
			if err != nil {
				log.Error("Failed to upload file to cloud service: " + err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to upload file to cloud service: "+err.Error()))
			}

			image := venue.VenuePictureCore{
				VenueID: venueID,
				URL:     fmt.Sprintf("%s%s", venueFileBaseURL, filepath.Base(filename)),
			}

			_, err = vh.service.CreateVenueImage(image)
			if err != nil {
				log.Error("Failed to insert image. " + err.Error())
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to insert image"))
			}
		}

		log.Sugar().Infof(venueID + " venue image added successfully")
		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "venue image added successfully"))
	}
}

func (vh *venueHandler) DeleteVenueImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := middlewares.ExtractToken(c)
		if err != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		venueImageID := c.Param("image_id")
		vn, err := vh.service.GetByVenueImageID(venueImageID)
		if err != nil {
			log.Error(err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
		}

		// Delete profile picture in the cloud before updated
		awsService := aws.InitS3()

		prevFilename := filepath.Base(vn.URL)
		prevPath := "profile-picture/" + prevFilename

		err = awsService.DeleteFile(prevPath)
		if err != nil {
			log.Error("Failed to delete file from cloud service: " + err.Error())
			return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Failed to delete file from cloud service: "+err.Error()))
		}

		err = vh.service.DeleteVenueImage(venueImageID)
		if err != nil {
			if strings.Contains(err.Error(), "failed to delete") {
				log.Error("failed to delete image")
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("failed to delete image"))
			} else {
				log.Error("Internal server error")
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("Internal server error"))
			}
		}

		log.Sugar().Infof(venueImageID + " venue image deleted successfully")
		return c.JSON(http.StatusOK, helper.SuccessResponse(nil, "venue image added successfully"))
	}
}

func (vh *venueHandler) GetAllVenueImage() echo.HandlerFunc {
	return func(c echo.Context) error {
		venueID := c.Param("venue_id")
		images, err := vh.service.GetAllVenueImage(venueID)
		if err != nil {
			log.Error(err.Error())
			return err
		}

		var resp []GetAllVenueImageResponse
		for _, image := range images {
			resp = append(resp, GetAllVenueImageToResponse(image))
		}

		return c.JSON(http.StatusOK, helper.SuccessResponse(resp, "venue image retrieved successfully"))
	}
}
