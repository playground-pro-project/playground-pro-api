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

		userId, err := middlewares.ExtractToken(c)
		if err != nil {
			log.Error("missing or malformed JWT")
			return c.JSON(http.StatusUnauthorized, helper.ResponseFormat(http.StatusUnauthorized, "Missing or Malformed JWT", nil, nil))
		}

		errBind := c.Bind(&req)
		if errBind != nil {
			log.Error("error on bind request")
			return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request"+errBind.Error(), nil, nil))
		}

		reservationData := requestReservation(req.Reservation)
		paymentData := req.Payment.requestPayment()
		reservation, payment, err := rh.service.MakeReservation(userId, reservationData, paymentData)
		if err != nil {
			switch {
			case strings.Contains(err.Error(), "empty"):
				log.Error("bad request, request cannot be empty")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request, request cannot be empty", nil, nil))
			case strings.Contains(err.Error(), "datetime"):
				log.Error("bad request, invalid datetime format")
				return c.JSON(http.StatusBadRequest, helper.ResponseFormat(http.StatusBadRequest, "Bad request, invalid datetime format", nil, nil))
			default:
				log.Error("internal server error")
				return c.JSON(http.StatusInternalServerError, helper.ResponseFormat(http.StatusInternalServerError, "Internal server error", nil, nil))
			}
		}

		response := makeReservation(payment)
		response.ReservationID = reservation.ReservationID
		return c.JSON(http.StatusCreated, helper.ResponseFormat(http.StatusCreated, "Successfully operation", response, nil))
	}
}
