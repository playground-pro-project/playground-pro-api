package service

import (
	"errors"
	"testing"
	"time"

	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/mocks"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
	"github.com/stretchr/testify/assert"
)

func TestCreateVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	userID := "useridstring"

	requestVenue := venue.VenueCore{
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

	requestVenuePictures := venue.VenuePictureCore{
		VenuePictureID: "venue_pictures_id_1",
		VenueID:        "venue_id_1",
		URL:            "https://contoh",
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

	t.Run("category cannot be empty", func(t *testing.T) {
		expectedErrorMessage := "category cannot be empty"
		requestVenue.Category = ""
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})

	t.Run("name cannot be empty", func(t *testing.T) {
		expectedErrorMessage := "name cannot be empty"
		requestVenue.Category = "category_1"
		requestVenue.Name = ""
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})

	t.Run("service time cannot be empty", func(t *testing.T) {
		expectedErrorMessage := "service time cannot be empty"
		requestVenue.Name = "venue_name_1"
		requestVenue.ServiceTime = ""
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})

	t.Run("location cannot be empty", func(t *testing.T) {
		expectedErrorMessage := "location cannot be empty"
		requestVenue.ServiceTime = "07:00 - 23:00"
		requestVenue.Location = ""
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})

	t.Run("price cannot be empty", func(t *testing.T) {
		expectedErrorMessage := "price cannot be empty"
		requestVenue.Location = "venue_location_1"
		requestVenue.Price = 0
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})

	t.Run("success create a venue", func(t *testing.T) {
		requestVenue.Price = 9.99
		data.On("InsertVenue", userID, requestVenue, requestVenuePictures).Return(expectedResult, nil).Once()
		result, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)
		data.AssertExpectations(t)
	})

	t.Run("error insert data, duplicated", func(t *testing.T) {
		expectedErrorMessage := "error insert data, duplicated"
		requestVenue.Price = 9.99
		data.On("InsertVenue", userID, requestVenue, requestVenuePictures).Return(venue.VenueCore{}, errors.New(expectedErrorMessage)).Once()
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		expectedErrorMessage := "internal server error"
		requestVenue.Price = 9.99
		data.On("InsertVenue", userID, requestVenue, requestVenuePictures).Return(venue.VenueCore{}, errors.New(expectedErrorMessage)).Once()
		_, err := service.CreateVenue(userID, requestVenue, requestVenuePictures)
		assert.NotNil(t, err)
		assert.EqualError(t, err, expectedErrorMessage, "Error message mismatch")
		data.AssertExpectations(t)
	})
}

func TestSearchVenues(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	keyword := "keyword"
	latitude := 123.456
	longitude := 789.012
	page := pagination.Pagination{
		Sort:       "",
		TotalRows:  10,
		TotalPages: 1,
	}

	t.Run("success", func(t *testing.T) {
		mockVenues := []venue.VenueCoreRaw{
			{
				VenueID:       "venue_id_1",
				OwnerID:       "owner_id_1",
				Category:      "category",
				Name:          "Venue 1",
				Description:   "Description 1",
				Username:      "username_1",
				ServiceTime:   "service_time",
				Location:      "location",
				Distance:      1.23,
				Price:         4.56,
				Longitude:     12.34,
				Latitude:      56.78,
				TotalRows:     100,
				TotalPages:    10,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
				DeletedAt:     time.Time{},
				TotalReviews:  20,
				AverageRating: 4.5,
				VenuePicture:  "venue_picture",
				VenuePictures: []venue.VenuePictureCore{},
				Reviews:       []venue.ReviewCore{},
				Reservations:  []venue.ReservationCore{},
				User:          venue.UserCore{},
			},
		}

		data.On("SearchVenues", keyword, latitude, longitude, page).Return(mockVenues, int64(10), 1, nil).Once()
		result, rows, pages, err := service.SearchVenues(keyword, latitude, longitude, page)
		assert.Nil(t, err)
		assert.Equal(t, mockVenues, result)
		assert.Equal(t, int64(10), rows)
		assert.Equal(t, 1, pages)
		data.AssertExpectations(t)
	})

	t.Run("venues not found", func(t *testing.T) {
		mockError := errors.New("venues not found")
		data.On("SearchVenues", keyword, latitude, longitude, page).Return([]venue.VenueCoreRaw{}, int64(0), 0, mockError).Once()
		result, rows, pages, err := service.SearchVenues(keyword, latitude, longitude, page)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "venues not found")
		assert.Equal(t, []venue.VenueCoreRaw{}, result)
		assert.Equal(t, int64(0), rows)
		assert.Equal(t, 0, pages)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("internal server error")
		data.On("SearchVenues", keyword, latitude, longitude, page).Return([]venue.VenueCoreRaw{}, int64(0), 0, mockError).Once()
		result, rows, pages, err := service.SearchVenues(keyword, latitude, longitude, page)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, []venue.VenueCoreRaw{}, result)
		assert.Equal(t, int64(0), rows)
		assert.Equal(t, 0, pages)
		data.AssertExpectations(t)
	})
}

func TestSelectVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	venueID := "venue_id_1"

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

	t.Run("success", func(t *testing.T) {
		data.On("SelectVenue", venueID).Return(expectedResult, nil).Once()
		result, err := service.SelectVenue(venueID)
		assert.Nil(t, err)
		assert.Equal(t, expectedResult, result)
		data.AssertExpectations(t)
	})

	t.Run("venue not found", func(t *testing.T) {
		data.On("SelectVenue", venueID).Return(venue.VenueCore{}, errors.New("not found, error while retrieving venue")).Once()
		result, err := service.SelectVenue(venueID)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, result)
		assert.ErrorContains(t, err, "not found, error while retrieving venue")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("SelectVenue", venueID).Return(venue.VenueCore{}, errors.New("internal server error")).Once()
		result, err := service.SelectVenue(venueID)
		assert.NotNil(t, err)
		assert.Equal(t, venue.VenueCore{}, result)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestEditVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	userID := "user_id_1"
	venueID := "venue_id_1"
	requestVenue := venue.VenueCore{
		Name:        "updated_venue_name",
		Description: "updated_venue_description",
	}

	t.Run("success", func(t *testing.T) {
		data.On("EditVenue", userID, venueID, requestVenue).Return(nil).Once()
		err := service.EditVenue(userID, venueID, requestVenue)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("venue profile record not found", func(t *testing.T) {
		data.On("EditVenue", userID, venueID, requestVenue).Return(errors.New("venue profile record not found")).Once()
		err := service.EditVenue(userID, venueID, requestVenue)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "venue profile record not found")
		data.AssertExpectations(t)
	})

	t.Run("no venue has been created", func(t *testing.T) {
		data.On("EditVenue", userID, venueID, requestVenue).Return(errors.New("no venue has been created")).Once()
		err := service.EditVenue(userID, venueID, requestVenue)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "no venue has been created")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("EditVenue", userID, venueID, requestVenue).Return(errors.New("internal server error")).Once()
		err := service.EditVenue(userID, venueID, requestVenue)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestUnregisterVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	userID := "user_id_1"
	venueID := "venue_id_1"

	t.Run("success", func(t *testing.T) {
		data.On("UnregisterVenue", userID, venueID).Return(nil).Once()
		err := service.UnregisterVenue(userID, venueID)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("venue record not found", func(t *testing.T) {
		data.On("UnregisterVenue", userID, venueID).Return(errors.New("venue record not found")).Once()
		err := service.UnregisterVenue(userID, venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "venue record not found")
		data.AssertExpectations(t)
	})

	t.Run("no row affected", func(t *testing.T) {
		data.On("UnregisterVenue", userID, venueID).Return(errors.New("no row affected")).Once()
		err := service.UnregisterVenue(userID, venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "no row affected")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("UnregisterVenue", userID, venueID).Return(errors.New("internal server error")).Once()
		err := service.UnregisterVenue(userID, venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		data.AssertExpectations(t)
	})
}

func TestVenueAvailability(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	venueID := "venue_id_1"

	t.Run("success", func(t *testing.T) {
		mockVenue := venue.VenueCore{
			VenueID:       venueID,
			Name:          "Venue 1",
			Description:   "Description 1",
			Location:      "Location 1",
			TotalRows:     1,
			TotalPages:    1,
			AverageRating: 4.5,
		}
		data.On("VenueAvailability", venueID).Return(mockVenue, nil).Once()
		result, err := service.VenueAvailability(venueID)
		assert.Nil(t, err)
		assert.Equal(t, mockVenue, result)
		data.AssertExpectations(t)
	})

	t.Run("venue not found", func(t *testing.T) {
		data.On("VenueAvailability", venueID).Return(venue.VenueCore{}, errors.New("not found")).Once()
		result, err := service.VenueAvailability(venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "list venues record not found")
		assert.Equal(t, venue.VenueCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("VenueAvailability", venueID).Return(venue.VenueCore{}, errors.New("internal server error")).Once()
		result, err := service.VenueAvailability(venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, venue.VenueCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestCreateVenueImage(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)

	t.Run("success", func(t *testing.T) {
		// Create a mock venue picture request
		mockReq := venue.VenuePictureCore{
			VenueID: "venue_id_1",
			URL:     "https://example.com/image.jpg",
		}

		// Mock the InsertVenueImage query method to return the mock venue picture
		data.On("InsertVenueImage", mockReq).Return(mockReq, nil).Once()

		// Call the CreateVenueImage method
		result, err := service.CreateVenueImage(mockReq)

		// Assert that there are no errors and the result matches the mock venue picture
		assert.Nil(t, err)
		assert.Equal(t, mockReq, result)

		// Assert that the mock data's method was called as expected
		data.AssertExpectations(t)
	})

	t.Run("missing venue ID", func(t *testing.T) {
		mockReq := venue.VenuePictureCore{
			VenueID: "",
			URL:     "https://example.com/image.jpg",
		}

		result, err := service.CreateVenueImage(mockReq)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error, venue ID is required")
		assert.Equal(t, venue.VenuePictureCore{}, result)
		data.AssertNotCalled(t, "InsertVenueImage")
	})

	t.Run("missing URL", func(t *testing.T) {
		mockReq := venue.VenuePictureCore{
			VenueID: "venue_id_1",
			URL:     "",
		}
		result, err := service.CreateVenueImage(mockReq)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "error, venue URL image is required")
		assert.Equal(t, venue.VenuePictureCore{}, result)
		data.AssertNotCalled(t, "InsertVenueImage")
	})

	t.Run("query error", func(t *testing.T) {
		mockReq := venue.VenuePictureCore{
			VenueID: "venue_id_1",
			URL:     "https://example.com/image.jpg",
		}

		mockError := errors.New("database error")
		data.On("InsertVenueImage", mockReq).Return(venue.VenuePictureCore{}, mockError).Once()
		result, err := service.CreateVenueImage(mockReq)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, mockError)
		assert.Equal(t, venue.VenuePictureCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestGetAllVenueImage(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	venueID := "venue_id_1"

	t.Run("success", func(t *testing.T) {
		// Create a mock list of venue images
		mockImages := []venue.VenuePictureCore{
			{VenuePictureID: "image_id_1", VenueID: venueID, URL: "https://example.com/image1.jpg"},
			{VenuePictureID: "image_id_2", VenueID: venueID, URL: "https://example.com/image2.jpg"},
		}
		data.On("GetAllVenueImage", venueID).Return(mockImages, nil).Once()
		result, err := service.GetAllVenueImage(venueID)
		assert.Nil(t, err)
		assert.Equal(t, mockImages, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("database error")
		data.On("GetAllVenueImage", venueID).Return([]venue.VenuePictureCore{}, mockError).Once()
		result, err := service.GetAllVenueImage(venueID)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, mockError)
		assert.Nil(t, result)
		data.AssertExpectations(t)
	})
}

func TestDeleteVenueImage(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	venueImageID := "image_id_1"

	t.Run("success", func(t *testing.T) {
		data.On("DeleteVenueImage", venueImageID).Return(nil).Once()
		err := service.DeleteVenueImage(venueImageID)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("database error")
		data.On("DeleteVenueImage", venueImageID).Return(mockError).Once()

		err := service.DeleteVenueImage(venueImageID)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, mockError)
		data.AssertExpectations(t)
	})
}

func TestGetVenueImageByID(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	venueID := "venue_id_1"
	venueImageID := "image_id_1"

	t.Run("success", func(t *testing.T) {
		// Create a mock venue image
		mockImage := venue.VenuePictureCore{
			VenuePictureID: venueImageID,
			VenueID:        venueID,
			URL:            "https://example.com/image1.jpg",
		}

		data.On("GetVenueImageByID", venueID, venueImageID).Return(mockImage, nil).Once()
		result, err := service.GetVenueImageByID(venueID, venueImageID)
		assert.Nil(t, err)
		assert.Equal(t, mockImage, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("database error")
		data.On("GetVenueImageByID", venueID, venueImageID).Return(venue.VenuePictureCore{}, mockError).Once()
		result, err := service.GetVenueImageByID(venueID, venueImageID)
		assert.NotNil(t, err)
		assert.ErrorIs(t, err, mockError)
		assert.Equal(t, venue.VenuePictureCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestMyVenues(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	userID := "user_id_1"

	t.Run("success", func(t *testing.T) {
		mockVenues := []venue.VenueCore{
			{VenueID: "venue_id_1", Name: "Venue 1"},
			{VenueID: "venue_id_2", Name: "Venue 2"},
		}

		data.On("MyVenues", userID).Return(mockVenues, nil).Once()
		result, err := service.MyVenues(userID)
		assert.Nil(t, err)
		assert.Equal(t, mockVenues, result)
		data.AssertExpectations(t)
	})

	t.Run("venues not found", func(t *testing.T) {
		mockError := errors.New("venues not found")
		data.On("MyVenues", userID).Return([]venue.VenueCore{}, mockError).Once()
		result, err := service.MyVenues(userID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "venues not found")
		assert.Equal(t, []venue.VenueCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("internal server error")
		data.On("MyVenues", userID).Return([]venue.VenueCore{}, mockError).Once()
		result, err := service.MyVenues(userID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, []venue.VenueCore{}, result)
		data.AssertExpectations(t)
	})

}
