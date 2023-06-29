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
	Name         string           `json:"venue_name,omitempty"`
	Location     string           `json:"location,omitempty"`
	CheckInDate  helper.LocalTime `json:"check_in_date,omitempty"`
	CheckOutDate helper.LocalTime `json:"check_out_date,omitempty"`
	Duration     float64          `json:"duration,omitempty"`
	Price        float64          `json:"price,omitempty"`
	GrandTotal   string           `json:"total_price,omitempty"`
	PaymentType  string           `json:"payment_type,omitempty"`
	PaymentCode  string           `json:"payment_code,omitempty"`
	Status       string           `json:"status,omitempty"`
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
