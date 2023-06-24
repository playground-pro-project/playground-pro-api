package service

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
)

var log = middlewares.Log()

type reservationService struct {
	query    reservation.ReservationData
	validate *validator.Validate
}

func New(rd reservation.ReservationData) reservation.ReservationService {
	return &reservationService{
		query:    rd,
		validate: validator.New(),
	}
}

// MakeReservation implements reservation.ReservationService.
func (rs *reservationService) MakeReservation(userId string, r reservation.ReservationCore, p reservation.PaymentCore) (reservation.ReservationCore, reservation.PaymentCore, error) {
	err := rs.validate.Struct(r)
	if err != nil {
		var message string
		switch {
		case strings.Contains(err.Error(), "venue_id"):
			log.Warn("venue_id cannot be empty")
			message = "venue_id cannot be empty"
		case strings.Contains(err.Error(), "format"):
			log.Warn("invalid datetime format")
			message = "invalid datetime format"
		case strings.Contains(err.Error(), "check_in_date"):
			log.Warn("check_in_date cannot be empty")
			message = "check_in_date cannot be empty"
		case strings.Contains(err.Error(), "check_out_date"):
			log.Warn("check_out_date cannot be empty")
			message = "check_out_date cannot be empty"
		default:
			log.Error("internal server error")
			message = "internal server error"
		}
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New(message)
	}

	result, paymentResult, err := rs.query.MakeReservation(userId, r, p)
	if err != nil {
		var message string
		switch {
		case strings.Contains(err.Error(), "user does not exist"):
			log.Error("failed to insert data, user does not exist")
			message = "user does not exist"
		case strings.Contains(err.Error(), "error insert data, duplicate entry"):
			log.Error("error insert data, duplicate entry")
			message = "error insert data, duplicate entry"
		default:
			log.Error("internal server error")
			message = "internal server error"
		}
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New(message)
	}

	log.Sugar().Infof("new reservation has been created: %s", result.ReservationID)
	return result, paymentResult, nil
}
