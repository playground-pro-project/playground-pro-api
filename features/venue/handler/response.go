package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

type searchVenueResponse struct {
	VenueID       string           `json:"venue_id,omitempty"`
	OwnerID       string           `json:"owner_id,omitempty"`
	Category      string           `json:"category,omitempty"`
	Name          string           `json:"name,omitempty"`
	Description   string           `json:"description,omitempty"`
	Username      string           `json:"username,omitempty"`
	ServiceTime   helper.LocalTime `json:"service_time,omitempty"`
	Location      string           `json:"location,omitempty"`
	Distance      uint             `json:"distance,omitempty"`
	Price         float64          `json:"price,omitempty"`
	TotalReviews  uint             `json:"total_reviews,omitempty"`
	AverageRating float64          `json:"average_rating,omitempty"`
	VenuePicture  string           `json:"venue_picture,omitempty"`
}

type VenuePicture struct {
	URL string
}

func searchVenue(v venue.VenueCore) searchVenueResponse {
	response := searchVenueResponse{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Username:      v.Username,
		ServiceTime:   helper.LocalTime(v.ServiceTime),
		Location:      v.Location,
		Distance:      v.Distance,
		Price:         v.Price,
		TotalReviews:  v.TotalReviews,
		AverageRating: v.AverageRating,
	}

	pictures := make([]VenuePicture, len(v.VenuePictures))
	for i, p := range v.VenuePictures {
		pictures[i] = VenuePicture{
			URL: p.URL,
		}
	}

	if len(pictures) > 0 {
		response.VenuePicture = pictures[0].URL
	}

	return response
}
