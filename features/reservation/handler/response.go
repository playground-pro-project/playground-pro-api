package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
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
