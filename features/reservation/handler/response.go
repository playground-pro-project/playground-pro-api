package handler

import (
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

type makeReservationResponse struct {
	PaymentID     string           `json:"payment_id"`
	ReservationID string           `json:"reservation_id"`
	GrandTotal    string           `json:"grand_total"`
	PaymentMethod string           `json:"payment_method"`
	PaymentType   string           `json:"payment_type"`
	PaymentCode   string           `json:"payment_code"`
	Status        string           `json:"status"`
	CreatedAt     helper.LocalTime `json:"created_at"`
	UpdatedAt     helper.LocalTime `json:"updated_at"`
}

func makeReservation(p reservation.PaymentCore) makeReservationResponse {
	return makeReservationResponse{
		PaymentID:     p.PaymentID,
		ReservationID: p.Reservation.ReservationID,
		GrandTotal:    p.GrandTotal,
		PaymentMethod: p.PaymentMethod,
		PaymentType:   p.PaymentType,
		PaymentCode:   p.PaymentCode,
		Status:        p.Status,
		CreatedAt:     helper.LocalTime(p.CreatedAt),
		UpdatedAt:     helper.LocalTime(p.UpdatedAt),
	}
}
