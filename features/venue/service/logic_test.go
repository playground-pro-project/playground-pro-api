package service

import (
	"errors"
	"testing"

	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/mocks"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	userID := "useridstring"
	request := venue.VenueCore{
		VenueID:     "venue_id_1",
		OwnerID:     "owner_id_1",
		Category:    "category_1",
		Name:        "venue_name_1",
		Description: "venue_description_1",
		Username:    "owner_username_1",
		ServiceTime: "07:00 - 23:00",
		Location:    "venue_location_1",
		Price:       9.99,
	}
	expectedResult := venue.VenueCore{
		VenueID:     "venue_id_1",
		OwnerID:     "owner_id_1",
		Category:    "category_1",
		Name:        "venue_name_1",
		Description: "venue_description_1",
		Username:    "owner_username_1",
		ServiceTime: "07:00 - 23:00",
		Location:    "venue_location_1",
		Price:       9.99,
	}
	emptyRequest := venue.VenueCore{
		VenueID:     "venue_id_1",
		OwnerID:     "owner_id_1",
		Category:    "",
		Name:        "",
		Description: "venue_description_1",
		Username:    "owner_username_1",
		ServiceTime: "",
		Location:    "",
		Price:       0,
	}

	t.Run("name cannot be empty", func(t *testing.T) {
		invalidRequest := emptyRequest
		invalidRequest.Name = "name_1"
		res, err := service.RegisterVenue(userID, invalidRequest)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, res)
		assert.ErrorContains(t, err, "name cannot be empty")
	})

	t.Run("service time cannot be empty", func(t *testing.T) {
		invalidRequest := emptyRequest
		invalidRequest.ServiceTime = "venue_name_1"
		res, err := service.RegisterVenue(userID, invalidRequest)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, res)
		assert.ErrorContains(t, err, "service time cannot be empty")
	})

	t.Run("location cannot be empty", func(t *testing.T) {
		invalidRequest := emptyRequest
		invalidRequest.Location = "07:00 - 23:00"
		res, err := service.RegisterVenue(userID, invalidRequest)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, res)
		assert.ErrorContains(t, err, "location cannot be empty")
	})

	t.Run("price cannot be empty", func(t *testing.T) {
		invalidRequest := emptyRequest
		invalidRequest.Location = "venue_location_1"
		res, err := service.RegisterVenue(userID, invalidRequest)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, res)
		assert.ErrorContains(t, err, "price cannot be empty")
	})

	t.Run("success create a venue", func(t *testing.T) {
		data.On("RegisterVenue", userID, request).Return(expectedResult, nil).Once()
		result, err := service.RegisterVenue(userID, request)
		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)
		data.AssertExpectations(t)
	})

	t.Run("error insert data, duplicated", func(t *testing.T) {
		data.On("RegisterVenue", mock.Anything, mock.Anything).Return(venue.VenueCore{}, errors.New("error insert data, duplicated")).Once()
		res, err := service.RegisterVenue(userID, request)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, res)
		assert.ErrorContains(t, err, "error insert data, duplicated")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("RegisterVenue", mock.Anything, mock.Anything).Return(venue.VenueCore{}, errors.New("internal server error")).Once()
		res, err := service.RegisterVenue(userID, request)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, res)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestSearchVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	keyword := "basket"
	latitude := 123.45
	longitude := 67.89

	page := pagination.Pagination{
		Limit:  3,
		Offset: 0,
		Page:   1,
	}
	expectedResult := []venue.VenueCore{
		{
			VenueID:       "venue_id_1",
			OwnerID:       "owner_id_1",
			Category:      "category_1",
			Name:          "venue_name_1",
			Description:   "venue_description_1",
			Username:      "owner_username_1",
			ServiceTime:   "07:00 - 23:00",
			Location:      "venue_location_1",
			Distance:      10,
			Price:         9.99,
			Longitude:     123.45,
			Latitude:      67.89,
			TotalReviews:  5,
			AverageRating: 4.5,
			VenuePictures: []venue.VenuePictureCore{
				{
					URL: "https://foto.com",
				},
			},
		},
	}
	expectedLimit := int(3)
	expectedPage := 1
	expectedTotalRows := int64(1)
	expectedTotalPages := 1

	t.Run("success", func(t *testing.T) {
		data.On("SearchVenues", keyword, page).Return(expectedResult, expectedTotalRows, expectedTotalPages, nil)
		result, totalRows, totalPages, err := service.SearchVenues(keyword, latitude, longitude, page)
		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedResult[0].VenueID, result[0].VenueID)
		assert.Equal(t, expectedResult[0].OwnerID, result[0].OwnerID)
		assert.Equal(t, expectedResult[0].Category, result[0].Category)
		// ... continue asserting other fields as needed
		assert.Equal(t, expectedLimit, page.Limit)
		assert.Equal(t, expectedPage, page.Page)
		assert.Equal(t, expectedTotalRows, totalRows)
		assert.Equal(t, expectedTotalPages, totalPages)
		data.AssertExpectations(t)
	})

	t.Run("list venues record not found", func(t *testing.T) {
		data.On("SearchVenues", keyword, page).Return([]venue.VenueCore{}, int64(0), 0, errors.New("list venues record not found")).Once()
		result, totalRows, totalPages, err := service.SearchVenues(keyword, latitude, longitude, page)
		assert.Error(t, err)
		assert.Len(t, result, 0)
		assert.EqualError(t, err, "list venues record not found")
		assert.Equal(t, expectedLimit, page.Limit)
		assert.Equal(t, expectedPage, page.Page)
		assert.Equal(t, int64(0), totalRows)
		assert.Equal(t, int64(0), totalPages)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("SearchVenues", keyword, page).Return([]venue.VenueCore{}, int64(0), 0, errors.New("internal server error"))
		result, totalRows, totalPages, err := service.SearchVenues(keyword, latitude, longitude, page)
		assert.Error(t, err)
		assert.Empty(t, result)
		assert.EqualError(t, err, "internal server error")
		assert.Equal(t, expectedLimit, page.Limit)
		assert.Equal(t, expectedPage, page.Page)
		assert.Equal(t, int64(0), totalRows)
		assert.Equal(t, int64(0), totalPages)
		data.AssertExpectations(t)
	})
}
