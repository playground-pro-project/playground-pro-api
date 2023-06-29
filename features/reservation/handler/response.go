package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

type makeReservationResponse struct {
	PaymentID     string `json:"payment_id"`
	ReservationID string `json:"reservation_id"`
	PaymentMethod string `json:"payment_method"`
	PaymentType   string `json:"payment_type"`
	PaymentCode   string `json:"payment_code"`
}

func makeReservation(p reservation.PaymentCore) makeReservationResponse {
	return makeReservationResponse{
		PaymentID:     p.PaymentID,
		ReservationID: p.Reservation.ReservationID,
		PaymentMethod: p.PaymentMethod,
		PaymentType:   p.PaymentType,
		PaymentCode:   p.PaymentCode,
	}
}

type reservationHistoryResponse struct {
	Name          string           `json:"venue_name,omitempty"`
	Location      string           `json:"location,omitempty"`
	CheckInDate   helper.LocalTime `json:"check_in_date,omitempty"`
	CheckOutDate  helper.LocalTime `json:"check_out_date,omitempty"`
	Duration      float64          `json:"duration,omitempty"`
	Price         float64          `json:"price,omitempty"`
	GrandTotal    string           `json:"total_price,omitempty"`
	PaymentType   string           `json:"payment_type,omitempty"`
	PaymentCode   string           `json:"payment_code,omitempty"`
	Status        string           `json:"status,omitempty"`
	ReservationID string           `json:"reservation_id,omitempty"`
	Venues        Venue            `json:"venues,omitempty"`
}

type Venue struct {
	Category     string                       `json:"category,omitempty"`
	Name         string                       `json:"venue_name,omitempty"`
	Description  string                       `json:"description,omitempty"`
	Username     string                       `json:"username,omitempty"`
	ServiceTime  string                       `json:"service_time,omitempty"`
	Price        float64                      `json:"price,omitempty"`
	Reservations []reservationHistoryResponse `json:"reservations,omitempty"`
}

func reservationHistory(payment reservation.PaymentCore) reservationHistoryResponse {
	response := reservationHistoryResponse{
		Name:         payment.Reservation.Venue.Name,
		Location:     payment.Reservation.Venue.Location,
		CheckInDate:  helper.LocalTime(payment.Reservation.CheckInDate),
		CheckOutDate: helper.LocalTime(payment.Reservation.CheckOutDate),
		Duration:     payment.Reservation.Duration,
		Price:        payment.Reservation.Venue.Price,
		GrandTotal:   payment.GrandTotal,
		PaymentType:  payment.PaymentType,
		PaymentCode:  payment.PaymentCode,
		Status:       payment.Status,
	}

	return response
}

type availability struct {
	ReservationID string           `json:"reservation_id,omitempty"`
	CheckInDate   helper.LocalTime `json:"check_in_date,omitempty"`
	CheckOutDate  helper.LocalTime `json:"check_out_date,omitempty"`
}

type venue struct {
	VenueID      string         `json:"venue_id,omitempty"`
	Category     string         `json:"category,omitempty"`
	Name         string         `json:"name,omitempty"`
	Reservations []availability `json:"reservations,omitempty"`
}

func Availability(reservations []reservation.AvailabilityCore) []venue {
	venuesMap := make(map[string]venue)
	for _, r := range reservations {
		venueID := r.VenueID
		v, ok := venuesMap[venueID]
		if !ok {
			v = venue{
				VenueID:      r.VenueID,
				Name:         r.Name,
				Category:     r.Category,
				Reservations: make([]availability, 0),
			}
		}

		reservation := availability{
			ReservationID: r.ReservationID,
			CheckInDate:   helper.LocalTime(r.CheckInDate),
			CheckOutDate:  helper.LocalTime(r.CheckOutDate),
		}

		v.Reservations = append(v.Reservations, reservation)
		venuesMap[venueID] = v
	}

	venues := make([]venue, 0, len(venuesMap))
	for _, v := range venuesMap {
		venues = append(venues, v)
	}

	return venues
}
