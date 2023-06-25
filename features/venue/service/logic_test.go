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

func TestSearchVenue(t *testing.T) {
	data := mocks.NewVenueData(t)
	service := New(data)
	keyword := "basket"
	page := pagination.Pagination{
		Limit:  3,
		Offset: 0,
	}
	expectedResult := []venue.VenueCore{
		{
			VenueID:       "venue_id_1",
			OwnerID:       "owner_id_1",
			Category:      "category_1",
			Name:          "venue_name_1",
			Description:   "venue_description_1",
			Username:      "owner_username_1",
			ServiceTime:   time.Now(),
			Location:      "venue_location_1",
			Distance:      10,
			Price:         9.99,
			Longitude:     123.45,
			Latitude:      67.89,
			TotalReviews:  5,
			AverageRating: 4.5,
			VenuePictures: []venue.VenuePictureCore{
				{
					VenuePictureID: "picture_id_1",
					VenueID:        "venue_id_1",
					URL:            "https://foto.com",
				},
			},
		},
	}
	expectedTotalRows := int64(1)
	expectedTotalPages := 1

	t.Run("success", func(t *testing.T) {
		data.On("SearchVenue", keyword, page).Return(expectedResult, expectedTotalRows, expectedTotalPages, nil)

		result, totalRows, totalPages, err := service.SearchVenues(keyword, page)

		assert.Nil(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, expectedResult[0].VenueID, result[0].VenueID)
		assert.Equal(t, expectedResult[0].OwnerID, result[0].OwnerID)
		assert.Equal(t, expectedResult[0].Category, result[0].Category)
		// ... continue asserting other fields as needed
		assert.Equal(t, expectedTotalRows, totalRows)
		assert.Equal(t, expectedTotalPages, totalPages)
		data.AssertExpectations(t)
	})

	t.Run("list venues record not found", func(t *testing.T) {
		data.On("SearchVenue", keyword, page).Return([]venue.VenueCore{}, int64(0), 0, errors.New("list venues record not found")).Once()

		result, _, _, err := service.SearchVenues(keyword, page)

		assert.Error(t, err)
		assert.Len(t, result, 0)
		assert.EqualError(t, err, "list venues record not found")
		data.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		data.On("SearchVenue", keyword, page).Return([]venue.VenueCore{}, int64(0), 0, errors.New("internal server error"))

		result, _, _, err := service.SearchVenues(keyword, page)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.EqualError(t, err, "internal server error")
		data.AssertExpectations(t)
	})

}
