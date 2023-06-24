package handler

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
)

type makeReservationRequest struct {
	VenueID      string `json:"venue_id" form:"venue_id"`
	CheckInDate  string `json:"check_in_date" form:"check_in_date" validate:"datetime"`
	CheckOutDate string `json:"check_out_date" form:"check_out_date" validate:"datetime"`
}

type createPaymentRequest struct {
	PaymentType string `json:"payment_type"  form:"payment_type"`
	GrandTotal  string `json:"grand_total"  form:"grand_total"`
}

type editReservationRequest struct {
	CheckInDate  *string `json:"check_in_date" form:"check_in_date" validate:"datetime"`
	CheckOutDate *string `json:"check_out_date" form:"check_out_date" validate:"datetime"`
}

func customDateTimeFormatValidator(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	_, err := time.Parse("2006-01-02 15:04:05", dateStr)
	return err == nil
}

func requestReservation(data interface{}) reservation.ReservationCore {
	result := reservation.ReservationCore{}
	validate := validator.New()
	validate.RegisterValidation("datetime", customDateTimeFormatValidator)

	err := validate.Struct(data)
	if err != nil {
		return result
	}

	switch v := data.(type) {
	case makeReservationRequest:
		result.VenueID = v.VenueID
		checkInDate, err := time.Parse("2006-01-02 15:04:05", v.CheckInDate)
		if err != nil {
			log.Error("error while parsing string to time format")
			return reservation.ReservationCore{}
		}
		result.CheckInDate = checkInDate
		checkOutDate, err := time.Parse("2006-01-02 15:04:05", v.CheckOutDate)
		if err != nil {
			log.Error("error while parsing string to time format")
			return reservation.ReservationCore{}
		}
		result.CheckOutDate = checkOutDate
	case *editReservationRequest:
		if v.CheckInDate != nil {
			checkInDate, err := time.Parse("2006-01-02 15:04:05", *v.CheckInDate)
			if err != nil {
				log.Error("error while parsing string to time format")
				return reservation.ReservationCore{}
			}
			result.CheckInDate = checkInDate
		}
		if v.CheckOutDate != nil {
			checkOutDate, err := time.Parse("2006-01-02 15:04:05", *v.CheckOutDate)
			if err != nil {
				log.Error("error while parsing string to time format")
				return reservation.ReservationCore{}
			}
			result.CheckOutDate = checkOutDate
		}
	default:
		return reservation.ReservationCore{}
	}

	return result
}

func (p createPaymentRequest) requestPayment() reservation.PaymentCore {
	return reservation.PaymentCore{
		PaymentType: p.PaymentType,
		GrandTotal:  p.GrandTotal,
	}
}
