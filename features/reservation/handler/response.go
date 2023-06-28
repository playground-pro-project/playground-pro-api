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
	Name         string           `json:"name,omitempty"`
	CheckInDate  helper.LocalTime `json:"check_in_date,omitempty"`
	CheckOutDate helper.LocalTime `json:"check_out_date,omitempty"`
	Duration     float64          `json:"duration,omitempty"`
	Price        float64          `json:"price,omitempty"`
	GrandTotal   string           `json:"total_price,omitempty"`
	PaymentType  string           `json:"payment_type,omitempty"`
	PaymentCode  string           `json:"payment_code,omitempty"`
	Status       string           `json:"status,omitempty"`
}

func reservationHistory(r reservation.ReservationCore) reservationHistoryResponse {
	response := reservationHistoryResponse{
		// Name:         r.Venues.Name,
		CheckInDate:  helper.LocalTime(r.CheckInDate),
		CheckOutDate: helper.LocalTime(r.CheckOutDate),
		Duration:     r.Duration,
		// Price:        r.Venues.Price,
		// GrandTotal:  r.Payment.GrandTotal,
		// PaymentType: r.Payment.PaymentType,
		// PaymentCode: r.Payment.PaymentCode,
		// Status:      r.Payment.Status,
	}

	return response
}
