package handler

import (
	"time"

	"github.com/playground-pro-project/playground-pro-api/features/venue"
)

type RegisterClassRequest struct {
	Category    string  `json:"category" form:"category"`
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	ServiceTime string  `json:"service_time" form:"service_time"`
	Location    string  `json:"location" form:"location"`
	Price       float64 `json:"price" form:"price"`
}

type EditClassRequest struct {
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
	case RegisterClassRequest:
		res.Category = v.Category
		res.Name = v.Name
		res.Description = v.Description
		serviceTime, err := time.ParseInLocation("07:45:00", v.ServiceTime, time.Local)
		if err != nil {
			log.Error("error while parsing string to time format")
			return venue.VenueCore{}
		}
		res.ServiceTime = serviceTime
		res.Location = v.Location
		res.Price = v.Price
	case *EditClassRequest:
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
			serviceTime, err := time.ParseInLocation("07:45:00", *v.ServiceTime, time.Local)
			if err != nil {
				log.Error("error while parsing string to time format")
				return venue.VenueCore{}
			}
			res.ServiceTime = serviceTime
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
