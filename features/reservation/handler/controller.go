package handler

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
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

// CheckAvailability implements reservation.ReservationHandler.
func (rh *reservationHandler) CheckAvailability() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		venueID := c.Param("venue_id")
		if venueID == "" {
			log.Error("empty venue_id parameter")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		availables, err := rh.service.CheckAvailability(venueID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				log.Info("there is no reservation at this moment")
				return helper.NotFoundError(c, "The requested resource was not found")
			}
			log.Error("internal server error")
			return helper.InternalServerError(c, "Internal server error")
		}

		if len(availables) == 0 {
			log.Error("there is no reservation yet atm")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		venues := Availability(availables)

		log.Sugar().Infoln(venues)
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successfully operation.", venues, nil))
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
			case strings.Contains(err.Error(), "reservation date not within the allowed timewindow"):
				log.Error("reservation date not within the allowed timewindow")
				return helper.BadRequestError(c, "Bad request")
			case strings.Contains(err.Error(), "reservation not available"):
				log.Error("reservation not available for the specified venue and timewindow")
				return helper.BadRequestError(c, "Bad request, reservation not available")
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

		res, err := rh.service.MyReservation(userId)
		if err != nil {
			if strings.Contains(err.Error(), "list reservations record not found") {
				log.Error("list reservations record not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		if len(res) == 0 {
			log.Error("reservation history not found")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		result := make([]myReservationResponse, len(res))
		for i, r := range res {
			result[i] = myReservation(r)
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
			if strings.Contains(err.Error(), "payment not found") {
				log.Error("payment not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		result, _ := reservationHistory(payment)
		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, nil))
	}
}

// MyVenueCharts implements reservation.ReservationHandler.
func (rh *reservationHandler) MyVenueCharts() echo.HandlerFunc {
	return func(c echo.Context) error {
		request := makeReservationRequest{}
		userId, errToken := middlewares.ExtractToken(c)
		if errToken != nil {
			log.Error("missing or malformed JWT")
			return helper.UnauthorizedError(c, "Missing or malformed JWT")
		}

		errBind := c.Bind(&request)
		if errBind != nil {
			log.Error("error on bind input")
			return helper.BadRequestError(c, "Bad request")
		}

		keyword := c.QueryParam("keyword")
		res, err := rh.service.MyVenueCharts(userId, keyword, reqReservation(request))
		if err != nil {
			if strings.Contains(err.Error(), "list charts record not found") {
				log.Error("list charts record not found")
				return helper.NotFoundError(c, "The requested resource was not found")
			} else {
				log.Error("internal server error")
				return helper.InternalServerError(c, "Internal server error")
			}
		}

		if len(res) == 0 {
			log.Error("list charts not found")
			return helper.NotFoundError(c, "The requested resource was not found")
		}

		result := make([]chartResponse, len(res))
		for i, r := range res {
			result[i] = charts(r)
		}

		return c.JSON(http.StatusOK, helper.ResponseFormat(http.StatusOK, "Successful Operation", result, nil))
	}
}
