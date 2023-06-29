package handler

import (
	"net/http"
	"strings"

	echo "github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
)

var log = middlewares.Log()

type reservationHandler struct {
	service reservation.ReservationService
}

func New(rs reservation.ReservationService) reservation.ReservationHandler {
	return &reservationHandler{
		service: rs,
	}
}

// MakeReservation implements reservation.ReservationHandler.
func (rh *reservationHandler) MakeReservation() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := struct {
			Reservation makeReservationRequest `json:"reservation"`
			Payment     createPaymentRequest   `json:"payment"`
		}{}

		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		errBind := c.Bind(&req)
		if errBind != nil {
			log.Error("error on bind request")
			return helper.BadRequestError(c, "Bad request")
		}

		reservationData := requestReservation(req.Reservation)
		paymentData := req.Payment.requestPayment()
		reservation, payment, err := rh.service.MakeReservation(userId, reservationData, paymentData)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "empty"):
				log.Error("bad request, request cannot be empty")
				return helper.BadRequestError(c, "Bad request")
			case strings.Contains(err.Error(), "datetime"):
				log.Error("bad request, invalid datetime format")
				return helper.BadRequestError(c, "Bad request")
			case strings.Contains(err.Error(), "unregistered user"):
				log.Error("unregistered user")
				return helper.BadRequestError(c, "Bad request")
			default:
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		response := makeReservation(payment)
		response.ReservationID = reservation.ReservationID
		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Successfully operation", response, nil))
	}
}

// ReservationStatus implements reservation.ReservationHandler.
func (rh *reservationHandler) ReservationStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		midtransResponse := midtransCallback{}
		log.Sugar().Info(midtransResponse)
		errBind := c.Bind(&midtransResponse)
		if errBind != nil {
			log.Sugar().Errorf("error on binding notification input", errBind)
			return helper.BadRequestError(c, "Bad request")
		}

		_, err := rh.service.ReservationStatus(reservationStatusRequest(midtransResponse))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Error("payment not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else if strings.Contains(err.Error(), "no payment record has been updated") {
				log.Error("no payment record has been updated")
				return helper.BadRequestError(c, "Bad request")
			} else if strings.Contains(err.Error(), "refund cannot be processed at least 1 hour away") {
				log.Error("refund cannot be processed at least 1 hour away")
				return helper.BadRequestError(c, "Bad request")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful updated payment status", nil, nil))
	}
}

// ReservationHistory implements reservation.ReservationHandler.
func (rh *reservationHandler) MyReservation() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		payments, err := rh.service.MyReservation(userId)
		if err != nil {
			if strings.Contains(err.Error(), "list reservations not found") {
				log.Error("reservations not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		if len(payments) == 0 {
			log.Error("reservation history not found")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		result := make([]reservationHistoryResponse, len(payments))
		for i, payment := range payments {
			result[i] = reservationHistory(payment)
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, nil))
	}
}

// DetailTransaction implements reservation.ReservationHandler.
func (rh *reservationHandler) DetailTransaction() echo.HandlerFunc {
	return func(c echo.Context) error {
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		paymentId := c.Param("payment_id")
		if paymentId == "" {
			log.Error("empty paymentId parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}
		payment, err := rh.service.DetailTransaction(userId, paymentId)
		if err != nil {
			if strings.Contains(err.Error(), "reservation not found") {
				log.Error("reservation not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		result := reservationHistory(payment)
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, nil))
	}
}
