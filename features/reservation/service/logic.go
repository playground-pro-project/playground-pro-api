package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	paymentgateway "github.com/playground-pro-project/playground-pro-api/utils/payment_gateway"
)

var log = middlewares.Log()

type reservationService struct {
	query    reservation.ReservationData
	validate *validator.Validate
	refund   paymentgateway.Refund
}

func New(rd reservation.ReservationData, refund paymentgateway.Refund) reservation.ReservationService {
	return &reservationService{
		query:    rd,
		validate: validator.New(),
		refund:   refund,
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

	// TODO 1 : Get price of spesific venue
	res1, err := rs.query.PriceVenue(r.VenueID)
	if err != nil {
		log.Sugar().Errorf("failed to get venue price %s", r.VenueID)
		return reservation.ReservationCore{}, reservation.PaymentCore{}, err
	}

	price, err := strconv.ParseFloat(fmt.Sprintf("%.2f", res1), 64)
	if err != nil {
		log.Sugar().Errorf("failed to parse grand_total: %s", err.Error())
		return reservation.ReservationCore{}, reservation.PaymentCore{}, err
	}

	log.Sugar().Infof("%.2f", price)

	// TODO 2: Get accumulative duration of specific venue
	duration := r.CheckOutDate.Sub(r.CheckInDate).Hours()
	r.Duration = duration
	log.Sugar().Infof("%.2f", r.Duration)

	// TODO 3: Multiply duration and price
	p.GrandTotal = strconv.FormatFloat(duration*price, 'f', 2, 64)

	log.Sugar().Infof(p.GrandTotal)

	result, paymentResult, err := rs.query.MakeReservation(userId, r, p)
	if err != nil {
		var message string
		switch {
		case strings.Contains(err.Error(), "user does not exist"):
			log.Error("failed to insert data, user does not exist")
			message = "user does not exist"
		default:
			log.Error("internal server error")
			message = "internal server error"
		}
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New(message)
	}

	log.Sugar().Infof("new reservation has been created: %s", result.ReservationID)
	return result, paymentResult, nil
}

// ReservationStatus implements reservation.ReservationService.
func (rs *reservationService) ReservationStatus(request reservation.PaymentCore) (reservation.PaymentCore, error) {
	switch request.Status {
	case "settlement":
		request.Status = "success"
		res, err := rs.query.ReservationStatus(request)
		if err != nil {
			log.Error("failed to update reservation status")
			return res, err
		}

		if !paymentgateway.IsRefundable(request.PaymentMethod) {
			grandTotal, errConv := strconv.ParseInt(request.GrandTotal, 10, 64)
			if errConv != nil {
				log.Error("failed to parse grand total")
				return res, errors.New("failed to parse grand total")
			}
			err := rs.refund.RefundTransaction(request.Reservation.ReservationID, grandTotal, "reason")
			if err != nil {
				log.Error("failed to refund transaction")
				return res, err
			}

			request.Status = "refund"
			res, err = rs.query.ReservationStatus(request)
			if err != nil {
				log.Error("failed to update refund status")
				return res, err
			}
		}

	case "expire":
		res, err := rs.query.ReservationStatus(request)
		if err != nil {
			log.Error("error on updating status to expire")
			return res, err
		}
	}

	return request, nil
}

// ReservationHistory implements reservation.ReservationService.
func (rs *reservationService) ReservationHistory(userId string) ([]reservation.PaymentCore, error) {
	payments, err := rs.query.ReservationHistory(userId)
	if err != nil {
		if strings.Contains(err.Error(), "list reservation not found") {
			log.Error("list reservation record not found")
			return []reservation.PaymentCore{}, errors.New("list reservation not found")
		} else {
			log.Error("internal server error")
			return []reservation.PaymentCore{}, errors.New("internal server error")
		}
	}

	return payments, err
}
