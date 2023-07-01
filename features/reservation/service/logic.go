package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

	// TODO 1 : Validate reservation timewindow and check availability
	minTime := time.Now().Local()
	maxTime := minTime.AddDate(0, 3, 0)
	if r.CheckInDate.Before(minTime) || r.CheckInDate.After(maxTime) {
		log.Warn("reservation date not within the allowed timewindow")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("reservation date not within the allowed timewindow")
	}

	if r.CheckOutDate.Before(minTime) || r.CheckOutDate.After(maxTime) {
		log.Warn("reservation date not within the allowed timewindow")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("reservation date not within the allowed timewindow")
	}

	// TODO 1.5: Check if there is an existing reservation for the same time slot
	existingReservations, err := rs.query.GetReservationsByTimeSlot(r.VenueID, r.CheckInDate, r.CheckOutDate)
	if err != nil {
		log.Sugar().Errorf("error on retrieving existing reservations: %s", err.Error())
		return reservation.ReservationCore{}, reservation.PaymentCore{}, err
	}

	if len(existingReservations) > 0 {
		log.Warn("reservation not available for the specified time slot")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("reservation not available")
	}

	// TODO 2 : Get price of spesific venue
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

	// TODO 3: Get accumulative duration of specific venue
	duration := r.CheckOutDate.Sub(r.CheckInDate).Hours()
	r.Duration = duration
	log.Sugar().Infof("%.2f", r.Duration)

	// TODO 4: Multiply duration and price
	p.GrandTotal = strconv.FormatFloat(duration*price, 'f', 2, 64)

	log.Sugar().Infof(p.GrandTotal)

	// TODO 5: Save data
	result, paymentResult, err := rs.query.MakeReservation(userId, r, p)
	if err != nil {
		var message string
		switch {
		case strings.Contains(err.Error(), "user does not exist"):
			log.Error("failed to insert data, user does not exist")
			message = "user does not exist"
		case strings.Contains(err.Error(), "unregistered user"):
			log.Error("foreign key constraint violation")
			message = "unregistered user"
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

	case "cancel":
		request.Status = "cancel"
		if !paymentgateway.IsRefundable(request.PaymentMethod) {
			grandTotal, errConv := strconv.ParseFloat(request.GrandTotal, 64)
			if errConv != nil {
				log.Error("failed to parse grand total")
				return request, errors.New("failed to parse grand total")
			}

			checkInTime, errQuery := rs.query.ReservationCheckOutDate(request.Reservation.ReservationID)
			if errQuery != nil {
				log.Sugar().Errorf("failed to get checkout date %v", checkInTime)
				return reservation.PaymentCore{}, errQuery
			}

			timeDiff := time.Until(checkInTime)
			if timeDiff < time.Hour {
				log.Error("refund cannot be processed at least 1 hour away")
				return request, errors.New("refund cannot be processed at least 1 hour away")
			}

			err := rs.refund.RefundTransaction(request.Reservation.ReservationID, int64(grandTotal), "reason")
			if err != nil {
				log.Error("failed to refund transaction")
				return request, err
			}
		}

		res, err := rs.query.ReservationStatus(request)
		if err != nil {
			log.Error("failed to update status to cancel")
			return res, err
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
func (rs *reservationService) MyReservation(userId string) ([]reservation.MyReservationCore, error) {
	payments, err := rs.query.MyReservation(userId)
	if err != nil {
		if strings.Contains(err.Error(), "list reservations record not found") {
			log.Error("list reservations record not found")
			return []reservation.MyReservationCore{}, errors.New("list reservations record not found")
		} else {
			log.Error("internal server error")
			return []reservation.MyReservationCore{}, errors.New("internal server error")
		}
	}

	return payments, err
}

func (rs *reservationService) DetailTransaction(userId string, paymentId string) (reservation.PaymentCore, error) {
	payment, err := rs.query.DetailTransaction(userId, paymentId)
	if err != nil {
		if strings.Contains(err.Error(), "payment not found") {
			log.Error("payment not found")
			return reservation.PaymentCore{}, errors.New("payment not found")
		} else {
			log.Error("internal server error")
			return reservation.PaymentCore{}, err
		}
	}

	return payment, nil
}

// CheckAvailability implements reservation.ReservationService.
func (rs *reservationService) CheckAvailability(venueId string) ([]reservation.AvailabilityCore, error) {
	result, err := rs.query.CheckAvailability(venueId)
	if err != nil {
		if strings.Contains(err.Error(), "list venues record not found") {
			log.Error("list venues record not found")
			return []reservation.AvailabilityCore{}, errors.New("list venues record not found")
		} else {
			log.Error("internal server error")
			return []reservation.AvailabilityCore{}, errors.New("internal server error")
		}
	}

	log.Sugar().Info(result)

	return result, err
}

// MyVenueCharts implements reservation.ReservationService.
func (rs *reservationService) MyVenueCharts(userId string, keyword string, request reservation.MyReservationCore) ([]reservation.MyReservationCore, error) {
	result, err := rs.query.MyVenueCharts(userId, keyword, request)
	if err != nil {
		if strings.Contains(err.Error(), "list charts record not found") {
			log.Error("list charts record not found")
			return []reservation.MyReservationCore{}, errors.New("list charts record not found")
		} else {
			log.Error("internal server error")
			return []reservation.MyReservationCore{}, errors.New("internal server error")
		}
	}

	return result, err
}
