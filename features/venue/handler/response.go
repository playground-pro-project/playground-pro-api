package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

type SearchVenueResponse struct {
	UserID        string  `json:"user_id,omitempty"`
	VenueID       string  `json:"venue_id,omitempty"`
	Category      string  `json:"category,omitempty"`
	Name          string  `json:"name,omitempty"`
	Username      string  `json:"username,omitempty"`
	Location      string  `json:"location,omitempty"`
	Distance      uint    `json:"distance,omitempty"`
	Price         float64 `json:"price,omitempty"`
	AverageRating float64 `json:"average_rating,omitempty"`
	VenuePicture  string  `json:"venue_picture,omitempty"`
}

type SelectVenueResponse struct {
	VenueID       string         `json:"venue_id,omitempty"`
	OwnerID       string         `json:"user_id,omitempty"`
	Category      string         `json:"category,omitempty"`
	Name          string         `json:"venue_name,omitempty"`
	Description   string         `json:"description,omitempty"`
	Username      string         `json:"username,omitempty"`
	ServiceTime   string         `json:"service_time,omitempty"`
	Location      string         `json:"location,omitempty"`
	Distance      uint           `json:"distance,omitempty"`
	Price         float64        `json:"price,omitempty"`
	TotalReviews  uint           `json:"total_reviews,omitempty"`
	AverageRating float64        `json:"average_rating,omitempty"`
	VenuePictures []VenuePicture `json:"venue_pictures,omitempty"`
	Reviews       []Review       `json:"reviews,omitempty"`
	Reservations  []Reservation  `json:"reservations,omitempty"`
}

type Review struct {
	Review string  `json:"review,omitempty"`
	Rating float64 `json:"rating,omitempty"`
}

type VenuePicture struct {
	VenuePictureURL string `json:"venue_picture_url,omitempty"`
}

type Reservation struct {
	ReservationID string           `json:"reservation_id,omitempty"`
	Username      string           `json:"username,omitempty"`
	CheckInDate   helper.LocalTime `json:"check_in_date,omitempty"`
	CheckOutDate  helper.LocalTime `json:"check_out_date,omitempty"`
}

type GetAllVenueImageResponse struct {
	VenuePictureID string `json:"venue_picture_id"`
	URL            string `json:"url"`
}

func GetAllVenueImageToResponse(v venue.VenuePictureCore) GetAllVenueImageResponse {
	return GetAllVenueImageResponse{
		VenuePictureID: v.VenuePictureID,
		URL:            v.URL,
	}
}

func SearchVenue(v venue.VenueCore) SearchVenueResponse {
	response := SearchVenueResponse{
		VenueID:       v.VenueID,
		Category:      v.Category,
		Name:          v.Name,
		Username:      v.Username,
		Location:      v.Location,
		Distance:      v.Distance,
		Price:         v.Price,
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
		Username:      v.Username,
		Description:   v.Description,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Distance:      v.Distance,
		Price:         v.Price,
		TotalReviews:  v.TotalReviews,
		AverageRating: v.AverageRating,
		VenuePictures: pictures,
		Reviews:       reviews,
	}

	return response
}

func Availability(a venue.VenueCore) SelectVenueResponse {
	reservations := make([]Reservation, len(a.Reservations))
	for i, r := range a.Reservations {
		reservations[i] = Reservation{
			ReservationID: r.ReservationID,
			Username:      r.Username,
			CheckInDate:   helper.LocalTime(r.CheckInDate),
			CheckOutDate:  helper.LocalTime(r.CheckOutDate),
		}
	}

	response := SelectVenueResponse{
		VenueID:      a.VenueID,
		OwnerID:      a.OwnerID,
		Category:     a.Category,
		Name:         a.Name,
		Reservations: reservations,
	}
	return response
}
