package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/venue"
)

type SearchVenueResponse struct {
	VenueID       string  `json:"venue_id,omitempty"`
	OwnerID       string  `json:"owner_id,omitempty"`
	Category      string  `json:"category,omitempty"`
	Name          string  `json:"name,omitempty"`
	Description   string  `json:"description,omitempty"`
	Username      string  `json:"username,omitempty"`
	ServiceTime   string  `json:"service_time,omitempty"`
	Location      string  `json:"location,omitempty"`
	Distance      uint    `json:"distance,omitempty"`
	Price         float64 `json:"price,omitempty"`
	TotalReviews  uint    `json:"total_reviews,omitempty"`
	AverageRating float64 `json:"average_rating,omitempty"`
	VenuePicture  string  `json:"venue_picture,omitempty"`
}

type SelectVenueResponse struct {
	VenueID       string         `json:"venue_id,omitempty"`
	OwnerID       string         `json:"owner_id,omitempty"`
	Category      string         `json:"category,omitempty"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Username      string         `json:"username,omitempty"`
	ServiceTime   string         `json:"service_time,omitempty"`
	Location      string         `json:"location,omitempty"`
	Distance      uint           `json:"distance,omitempty"`
	Price         float64        `json:"price,omitempty"`
	TotalReviews  uint           `json:"total_reviews,omitempty"`
	AverageRating float64        `json:"average_rating,omitempty"`
	VenuePictures []VenuePicture `json:"venue_picture,omitempty"`
	Reviews       []Review       `json:"reviews,omitempty"`
}

type Review struct {
	Review string  `json:"review" form:"review"`
	Rating float64 `json:"rating" form:"rating"`
}

type VenuePicture struct {
	VenuePictureURL string `json:"venue_picture_url" form:"venue_picture_url"`
}

func SearchVenue(v venue.VenueCore) SearchVenueResponse {
	response := SearchVenueResponse{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Username:      v.Username,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Distance:      v.Distance,
		Price:         v.Price,
		TotalReviews:  v.TotalReviews,
		AverageRating: v.AverageRating,
	}

	pictures := make([]VenuePicture, len(v.VenuePictures))
	for i, p := range v.VenuePictures {
		pictures[i] = VenuePicture{
			VenuePictureURL: p.URL,
		}
	}

	if len(pictures) > 0 {
		response.VenuePicture = pictures[0].VenuePictureURL
	}

	return response
}

func SelectVenue(v venue.VenueCore) SelectVenueResponse {
	pictures := make([]VenuePicture, len(v.VenuePictures))
	for i, p := range v.VenuePictures {
		pictures[i] = VenuePicture{
			VenuePictureURL: p.URL,
		}
	}

	reviews := make([]Review, len(v.Reviews))
	for i, r := range v.Reviews {
		reviews[i] = Review{
			Review: r.Review,
			Rating: r.Rating,
		}
	}

	response := SelectVenueResponse{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Username:      v.Username,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Distance:      v.Distance,
		Price:         v.Price,
		TotalReviews:  v.TotalReviews,
		AverageRating: v.AverageRating,
	}

	return response
}
