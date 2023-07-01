package handler

import (
	"fmt"

	"github.com/playground-pro-project/playground-pro-api/features/venue"
)

const (
	maxVenueFileSize = 2 * 1 << 20 // 2 MB
	venueFileBaseURL = "https://aws-pgp-bucket.s3.ap-southeast-2.amazonaws.com/venue-images/"
)

type RegisterVenueRequest struct {
	Category    string  `json:"category" form:"category"`
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	ServiceTime string  `json:"service_time" form:"service_time"`
	Location    string  `json:"location" form:"location"`
	Price       float64 `json:"price" form:"price"`
	Longitude   float64 `json:"lon" form:"lon"`
	Latitude    float64 `json:"lat" form:"lat"`
}

type EditVenueRequest struct {
	Category    *string  `json:"category" form:"category"`
	Name        *string  `json:"name" form:"name"`
	Description *string  `json:"description" form:"description"`
	ServiceTime *string  `json:"service_time" form:"service_time"`
	Location    *string  `json:"location" form:"location"`
	Price       *float64 `json:"price" form:"price"`
}

func RequestToCore(data interface{}) venue.VenueCore {
	res := venue.VenueCore{}
	switch v := data.(type) {
	case RegisterVenueRequest:
		res.Category = v.Category
		res.Name = v.Name
		res.Description = v.Description
		res.ServiceTime = v.ServiceTime
		res.Location = v.Location
		res.Price = v.Price
		res.Longitude = v.Longitude
		res.Latitude = v.Latitude
	case *EditVenueRequest:
		if v.Category != nil {
			res.Category = *v.Category
		}
		if v.Name != nil {
			res.Name = *v.Name
		}
		if v.Description != nil {
			res.Description = *v.Description
		}
		if v.ServiceTime != nil {
			res.ServiceTime = *v.ServiceTime
		}
		if v.Location != nil {
			res.Location = *v.Location
		}
		if v.Price != nil {
			res.Price = *v.Price
		}
	default:
		return venue.VenueCore{}

	}
	return res
}

func validateRegisterVenueRequest(request RegisterVenueRequest) error {
	if request.Category == "" {
		return fmt.Errorf("category is required")
	}
	if request.Name == "" {
		return fmt.Errorf("name is required")
	}
	if request.Description == "" {
		return fmt.Errorf("description is required")
	}
	if request.ServiceTime == "" {
		return fmt.Errorf("service time is required")
	}
	if request.Location == "" {
		return fmt.Errorf("location is required")
	}
	if request.Price == 0 {
		return fmt.Errorf("price is required")
	}
	if request.Longitude == 0 {
		return fmt.Errorf("longitude is required")
	}
	if request.Latitude == 0 {
		return fmt.Errorf("latitude is required")
	}

	return nil
}
