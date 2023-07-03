package service

import (
	"errors"
	"testing"
	"time"

	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	"github.com/playground-pro-project/playground-pro-api/mocks"
	paymentgateway "github.com/playground-pro-project/playground-pro-api/utils/payment_gateway"
	"github.com/stretchr/testify/assert"
)

func TestMyVenueCharts(t *testing.T) {
	data := mocks.NewReservationData(t)
	refund := &paymentgateway.MyRefund{}
	service := New(data, refund)
	userID := "user_id_1"
	keyword := "keyword"
	checkInDateStr := "2022-12-25 20:00:00"
	checkOutDateStr := "2022-12-25 21:00:00"
	checkInDate, _ := time.Parse("2006-01-02", checkInDateStr)
	checkOutDate, _ := time.Parse("2006-01-02", checkOutDateStr)
	request := reservation.MyReservationCore{
		// Set the necessary fields for the request object
	}

	t.Run("success", func(t *testing.T) {
		mockReservations := []reservation.MyReservationCore{
			{VenueID: "venue_id_1", VenueName: "Venue 1"},
			{VenueID: "venue_id_2", VenueName: "Venue 2"},
		}

		data.On("MyVenueCharts", userID, keyword, request).Return(mockReservations, nil).Once()
		result, err := service.MyVenueCharts(userID, keyword, checkInDate, checkOutDate)
		assert.Nil(t, err)
		assert.Equal(t, mockReservations, result)
		data.AssertExpectations(t)
	})

	t.Run("list charts record not found", func(t *testing.T) {
		mockError := errors.New("list charts record not found")
		data.On("MyVenueCharts", userID, keyword, request).Return([]reservation.MyReservationCore{}, mockError).Once()
		result, err := service.MyVenueCharts(userID, keyword, checkInDate, checkOutDate)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "list charts record not found")
		assert.Equal(t, []reservation.MyReservationCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("internal server error")
		data.On("MyVenueCharts", userID, keyword, request).Return([]reservation.MyReservationCore{}, mockError).Once()
		result, err := service.MyVenueCharts(userID, keyword, checkInDate, checkOutDate)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, []reservation.MyReservationCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestCheckAvailability(t *testing.T) {
	data := mocks.NewReservationData(t)
	refund := &paymentgateway.MyRefund{}
	service := New(data, refund)
	venueID := "venue_id_1"

	t.Run("success", func(t *testing.T) {
		mockAvailability := []reservation.AvailabilityCore{
			{
				VenueID:       "venue_id_1",
				Name:          "Venue 1",
				Category:      "Category 1",
				PaymentID:     "payment_id_1",
				ReservationID: "reservation_id_1",
				CheckInDate:   time.Now(),
				CheckOutDate:  time.Now().AddDate(0, 0, 1),
			},
			{
				VenueID:       "venue_id_2",
				Name:          "Venue 2",
				Category:      "Category 2",
				PaymentID:     "payment_id_2",
				ReservationID: "reservation_id_2",
				CheckInDate:   time.Now(),
				CheckOutDate:  time.Now().AddDate(0, 0, 2),
			},
		}

		data.On("CheckAvailability", venueID).Return(mockAvailability, nil).Once()
		result, err := service.CheckAvailability(venueID)
		assert.Nil(t, err)
		assert.Equal(t, mockAvailability, result)
		data.AssertExpectations(t)
	})

	t.Run("list venues record not found", func(t *testing.T) {
		mockError := errors.New("list venues record not found")
		data.On("CheckAvailability", venueID).Return([]reservation.AvailabilityCore{}, mockError).Once()
		result, err := service.CheckAvailability(venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "list venues record not found")
		assert.Equal(t, []reservation.AvailabilityCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("internal server error")
		data.On("CheckAvailability", venueID).Return([]reservation.AvailabilityCore{}, mockError).Once()
		result, err := service.CheckAvailability(venueID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, []reservation.AvailabilityCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestDetailTransaction(t *testing.T) {
	data := mocks.NewReservationData(t)
	refund := &paymentgateway.MyRefund{}
	service := New(data, refund)
	userID := "user_id_1"
	paymentID := "payment_id_1"

	t.Run("success", func(t *testing.T) {
		mockPayment := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			PaymentType:   "type",
			PaymentCode:   "code",
			GrandTotal:    "100",
			ServiceFee:    10.0,
			Status:        "paid",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			ReservationID: "reservation_id_1",
			Reservation:   reservation.ReservationCore{},
			Venue:         reservation.VenueCore{},
		}

		data.On("DetailTransaction", userID, paymentID).Return(mockPayment, nil).Once()
		result, err := service.DetailTransaction(userID, paymentID)
		assert.Nil(t, err)
		assert.Equal(t, mockPayment, result)
		data.AssertExpectations(t)
	})

	t.Run("payment not found", func(t *testing.T) {
		mockError := errors.New("payment not found")
		data.On("DetailTransaction", userID, paymentID).Return(reservation.PaymentCore{}, mockError).Once()
		result, err := service.DetailTransaction(userID, paymentID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "payment not found")
		assert.Equal(t, reservation.PaymentCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("internal server error")
		data.On("DetailTransaction", userID, paymentID).Return(reservation.PaymentCore{}, mockError).Once()
		result, err := service.DetailTransaction(userID, paymentID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, reservation.PaymentCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestMyReservation(t *testing.T) {
	data := mocks.NewReservationData(t)
	refund := &paymentgateway.MyRefund{}
	service := New(data, refund)
	userID := "user_id_1"

	t.Run("success", func(t *testing.T) {
		mockReservations := []reservation.MyReservationCore{
			{
				VenueID:       "venue_id_1",
				VenueName:     "Venue 1",
				Location:      "Location 1",
				ReservationID: "reservation_id_1",
				CheckInDate:   time.Now(),
				CheckOutDate:  time.Now(),
				Duration:      1.5,
				Price:         100.0,
				PaymentID:     "payment_id_1",
				GrandTotal:    150.0,
				PaymentType:   "type",
				PaymentCode:   "code",
				Status:        "paid",
				SalesVolume:   10,
			},
		}

		data.On("MyReservation", userID).Return(mockReservations, nil).Once()
		result, err := service.MyReservation(userID)
		assert.Nil(t, err)
		assert.Equal(t, mockReservations, result)
		data.AssertExpectations(t)
	})

	t.Run("list reservations record not found", func(t *testing.T) {
		mockError := errors.New("list reservations record not found")
		data.On("MyReservation", userID).Return([]reservation.MyReservationCore{}, mockError).Once()
		result, err := service.MyReservation(userID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "list reservations record not found")
		assert.Equal(t, []reservation.MyReservationCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("query error", func(t *testing.T) {
		mockError := errors.New("internal server error")
		data.On("MyReservation", userID).Return([]reservation.MyReservationCore{}, mockError).Once()
		result, err := service.MyReservation(userID)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "internal server error")
		assert.Equal(t, []reservation.MyReservationCore{}, result)
		data.AssertExpectations(t)
	})
}

func TestReservationStatus(t *testing.T) {
	data := mocks.NewReservationData(t)
	refund := &mocks.Refund{}
	service := New(data, refund)

	t.Run("success - expire", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "expire",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}

		data.On("ReservationStatus", request).Return(request, nil).Once()

		result, err := service.ReservationStatus(request)
		assert.Nil(t, err)
		assert.Equal(t, "expire", result.Status)
		data.AssertExpectations(t)
	})

	t.Run("error - failed to parse grand total", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "cancel",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "invalid",
		}

		result, err := service.ReservationStatus(request)
		assert.Error(t, err)
		assert.Equal(t, "failed to parse grand total: strconv.ParseFloat: parsing \"invalid\": invalid syntax", err.Error())
		assert.Equal(t, request, result)
	})

	t.Run("error - failed to get checkout date", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "cancel",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}

		data.On("ReservationCheckOutDate", request.Reservation.ReservationID).Return(time.Time{}, errors.New("database error")).Once()

		result, err := service.ReservationStatus(request)
		assert.Error(t, err)
		assert.Equal(t, "failed to get checkout date: database error", err.Error())
		assert.Equal(t, reservation.PaymentCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("error - refund cannot be processed at least 1 hour away", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "cancel",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}

		mockCheckOutDate := time.Now().Add(time.Minute * 30)
		data.On("ReservationCheckOutDate", request.Reservation.ReservationID).Return(mockCheckOutDate, nil).Once()

		result, err := service.ReservationStatus(request)
		assert.Error(t, err)
		assert.Equal(t, "refund cannot be processed at least 1 hour away", err.Error())
		assert.Equal(t, request, result)
		data.AssertExpectations(t)
	})

	t.Run("error - failed to refund transaction", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "cancel",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}

		mockCheckOutDate := time.Now().Add(time.Hour * 2)
		data.On("ReservationCheckOutDate", request.Reservation.ReservationID).Return(mockCheckOutDate, nil).Once()
		refund.On("RefundTransaction", request.Reservation.ReservationID, int64(150.0), "reason").Return(errors.New("refund error")).Once()

		result, err := service.ReservationStatus(request)
		assert.Error(t, err)
		assert.Equal(t, "failed to refund transaction: refund error", err.Error())
		assert.Equal(t, request, result)
		data.AssertExpectations(t)
		refund.AssertExpectations(t)
	})

	t.Run("error - error on updating status to expire", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "expire",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}

		data.On("ReservationStatus", request).Return(reservation.PaymentCore{}, errors.New("database error")).Once()
		result, err := service.ReservationStatus(request)
		assert.Error(t, err)
		assert.Equal(t, "error on updating status to expire: database error", err.Error())
		assert.Equal(t, reservation.PaymentCore{}, result)
		data.AssertExpectations(t)
	})

	t.Run("success - settlement", func(t *testing.T) {
		request := reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "settlement",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}
		data.On("ReservationStatus", request).Return(reservation.PaymentCore{
			PaymentID:     "payment_id_1",
			PaymentMethod: "method",
			Status:        "success",
			Reservation: reservation.ReservationCore{
				ReservationID: "reservation_id_1",
			},
			GrandTotal: "150.0",
		}, nil).Once()
		result, err := service.ReservationStatus(request)
		assert.Nil(t, err)
		assert.Equal(t, "success", result.Status)
		data.AssertExpectations(t)
	})
}

func TestMakeReservation(t *testing.T) {
	data := mocks.NewReservationData(t)
	refund := &paymentgateway.MyRefund{}
	service := New(data, refund)
	userId := "user_id_1"

	t.Run("success", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{}, nil).Once()
		data.On("PriceVenue", reservationCore.VenueID).Return(100.0, nil).Once()
		data.On("MakeReservation", userId, reservationCore, paymentCore).Return(reservationCore, paymentCore, nil).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)
		assert.Nil(t, err)
		assert.Equal(t, reservationCore, result)
		assert.Equal(t, paymentCore, paymentResult)
		data.AssertExpectations(t)
	})

	t.Run("error - venue_id is empty", func(t *testing.T) {
		userId := "user_id_1"
		reservationCore := reservation.ReservationCore{
			VenueID:      "",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "venue_id cannot be empty", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
	})

	t.Run("error - check_in_date is empty", func(t *testing.T) {
		userId := "user_id_1"
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Time{},
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "check_in_date cannot be empty", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
	})

	t.Run("error - check_out_date is empty", func(t *testing.T) {
		userId := "user_id_1"
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Time{},
		}
		paymentCore := reservation.PaymentCore{}

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "check_out_date cannot be empty", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
	})

	t.Run("error - reservation date not within the allowed timewindow", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, -1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
	})

	t.Run("error - reservation not available for the specified time slot", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods to simulate an existing reservation
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{{}}, nil).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "reservation not available", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
		data.AssertExpectations(t)
	})

	t.Run("error - failed to get venue price", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods to simulate an error in getting the venue price
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{}, nil).Once()
		data.On("PriceVenue", reservationCore.VenueID).Return(0.0, errors.New("failed to get venue price")).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "failed to get venue price", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
		data.AssertExpectations(t)
	})

	t.Run("error - failed to parse grand_total", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods to simulate an error in parsing grand_total
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{}, nil).Once()
		data.On("PriceVenue", reservationCore.VenueID).Return(100.0, nil).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "failed to parse grand_total", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
		data.AssertExpectations(t)
	})

	t.Run("error - failed to insert data, user does not exist", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods to simulate a user-related error
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{}, nil).Once()
		data.On("PriceVenue", reservationCore.VenueID).Return(100.0, nil).Once()
		data.On("MakeReservation", userId, reservationCore, paymentCore).Return(reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("user does not exist")).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "user does not exist", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
		data.AssertExpectations(t)
	})

	t.Run("error - foreign key constraint violation", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods to simulate a foreign key constraint violation error
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{}, nil).Once()
		data.On("PriceVenue", reservationCore.VenueID).Return(100.0, nil).Once()
		data.On("MakeReservation", userId, reservationCore, paymentCore).Return(reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("unregistered user")).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "unregistered user", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
		data.AssertExpectations(t)
	})

	t.Run("error - internal server error", func(t *testing.T) {
		reservationCore := reservation.ReservationCore{
			VenueID:      "venue_id_1",
			CheckInDate:  time.Now().AddDate(0, 0, 1),
			CheckOutDate: time.Now().AddDate(0, 0, 2),
		}
		paymentCore := reservation.PaymentCore{}

		// Mock the necessary methods to simulate an internal server error
		data.On("GetReservationsByTimeSlot", reservationCore.VenueID, reservationCore.CheckInDate, reservationCore.CheckOutDate).Return([]reservation.ReservationCore{}, nil).Once()
		data.On("PriceVenue", reservationCore.VenueID).Return(100.0, nil).Once()
		data.On("MakeReservation", userId, reservationCore, paymentCore).Return(reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error")).Once()

		result, paymentResult, err := service.MakeReservation(userId, reservationCore, paymentCore)

		assert.Error(t, err)
		assert.Equal(t, "internal server error", err.Error())
		assert.Equal(t, reservation.ReservationCore{}, result)
		assert.Equal(t, reservation.PaymentCore{}, paymentResult)
		data.AssertExpectations(t)
	})
}
