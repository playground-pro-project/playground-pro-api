package data

import (
	"errors"
	"strings"
	"time"

	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
	paymentgateway "github.com/playground-pro-project/playground-pro-api/utils/payment_gateway"
	"gorm.io/gorm"
)

var log = middlewares.Log()

type reservationQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) reservation.ReservationData {
	return &reservationQuery{
		db: db,
	}
}

// MakeReservation implements reservation.ReservationData.
func (rq *reservationQuery) MakeReservation(userID string, r reservation.ReservationCore, p reservation.PaymentCore) (reservation.ReservationCore, reservation.PaymentCore, error) {
	tx := rq.db.Begin()
	if tx.Error != nil {
		log.Error("error on beginning database transaction")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error on beginning database transaction")
	}

	r.UserID = userID
	reservationModel := reservationEntities(r)
	reservationModel.ReservationID = helper.GenerateReservationID()

	// TODO 1 : Create reservation
	if err := tx.Create(&reservationModel).Error; err != nil {
		tx.Rollback()
		log.Error("error while creating reservation")
		if strings.Contains(err.Error(), "Error 1452") {
			return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("unregistered user")
		}
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error while creating reservation")
	}

	log.Sugar().Info(reservationModel)

	// TODO 2 : Charge payment using Midtrans
	paymentModel, err := paymentgateway.ChargeMidtrans(reservationModel.ReservationID, p)
	if err != nil {
		log.Error("error while charging Midtrans payment")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error while charging Midtrans payment")
	}

	// TODO 3 : Create payment
	if err := tx.Create(paymentEntities(PaymentCoreFromChargeResponse(paymentModel))).Error; err != nil {
		tx.Rollback()
		log.Error("error while saving payment")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error while saving payment")
	}

	// TODO 4 : Assign payment ID to reservation
	reservationModel.PaymentID = &paymentModel.TransactionID
	if err := tx.Save(&reservationModel).Error; err != nil {
		tx.Rollback()
		log.Error("error while updating reservation with payment_id")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error while updating reservation with payment_id")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Error("error on committing database transaction")
		return reservation.ReservationCore{}, reservation.PaymentCore{}, errors.New("internal server error on committing database transaction")
	}

	return reservationModels(reservationModel), PaymentCoreFromChargeResponse(paymentModel), nil
}

// TODO 5: Callback Midtrans for updated payment status during reservation validation
// IF status: settlement then success, ELIF status: expired/cancel(refund) then failed.
// ReservationStatus implements reservation.ReservationData.
func (rq *reservationQuery) ReservationStatus(request reservation.PaymentCore) (reservation.PaymentCore, error) {
	req := paymentEntities(request)
	query := rq.db.Table("payments").
		Where("payment_id = ?", request.PaymentID).
		Updates(map[string]interface{}{
			"status": request.Status,
		})
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("user profile record not found")
		return reservation.PaymentCore{}, errors.New("user profile record not found")
	}

	if query.RowsAffected == 0 {
		log.Warn("no payment record has been updated")
		return reservation.PaymentCore{}, errors.New("no row affected")
	}

	if query.Error != nil {
		log.Error("error while updating payment status")
		return reservation.PaymentCore{}, errors.New("internal server error")
	}

	return paymentModels(*req), nil
}

// PriceVenue retrieves the price of a venue by its ID
func (rq *reservationQuery) PriceVenue(venue_id string) (float64, error) {
	venue := Venue{}
	query := rq.db.Table("venues").
		Select("venues.price").
		Where("venue_id = ?", venue_id).
		First(&venue)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue not found")
		return 0, errors.New("venue not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing venue query:", query.Error)
		return 0, query.Error
	}
	log.Sugar().Infof("venue data found in the database %f", venue.Price)
	return venue.Price, nil
}

// ReservationCheckOutDate retrieves check out time of a reservation by its ID
func (rq *reservationQuery) ReservationCheckOutDate(reservation_id string) (time.Time, error) {
	reservation := Reservation{}
	query := rq.db.Table("reservations").
		Select("reservations.check_out_date").
		Where("reservation_id = ?", reservation_id).
		First(&reservation)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue not found")
		return time.Time{}, errors.New("venue not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing venue query:", query.Error)
		return time.Time{}, query.Error
	}
	log.Sugar().Infof("checkout date found in the database %v", reservation.CheckOutDate)
	return reservation.CheckOutDate, nil
}

// ReservationHistory implements reservation.ReservationData.
func (rq *reservationQuery) MyReservation(userId string) ([]reservation.MyReservationCore, error) {
	result := []MyReservation{}
	query := rq.db.Raw(`
		SELECT venues.venue_id,
			venues.name,
			venues.location,
			venues.price,
			reservations.reservation_id, 
			reservations.check_in_date,
			reservations.check_out_date,	
			reservations.duration,
			payments.payment_id,
			payments.payment_type,
			payments.payment_code,
			payments.status
		FROM payments
		INNER JOIN reservations ON payments.payment_id = reservations.payment_id
		INNER JOIN venues ON reservations.venue_id = venues.venue_id
		WHERE reservations.user_id = ?
	`, userId).Scan(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list reservations record not found")
		return nil, errors.New("list reservations record not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing list reservations query:", query.Error)
		return nil, query.Error
	} else {
		log.Sugar().Info("list reservations data found in the database")
	}

	return modelToMyReservationCore(result), nil
}

// GetVenueNameAndPrice retrieves the name and price of a venue by its ID
func (rq *reservationQuery) GetVenueNameAndPrice(venueID string) (string, string, float64, error) {
	venue := Venue{}
	query := rq.db.Raw("SELECT name, location, price FROM venues WHERE venue_id = ?", venueID).Scan(&venue)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue not found")
		return "", "", 0, errors.New("venue not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing venue query:", query.Error)
		return "", "", 0, query.Error
	}
	log.Sugar().Infof("venue data found in the database: Name=%s, Price=%f", venue.Name, venue.Price)
	return venue.Name, venue.Location, venue.Price, nil
}

// DetailTransaction implements reservation.ReservationData.
func (rq *reservationQuery) DetailTransaction(userId string, paymentId string) (reservation.PaymentCore, error) {
	payment := Payment{}
	query := rq.db.Preload("Reservation").Where("payment_id = ?", paymentId).First(&payment)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("payment not found")
		return reservation.PaymentCore{}, errors.New("payment not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing payment query:", query.Error)
		return reservation.PaymentCore{}, query.Error
	}

	if payment.Reservation.VenueID != "" {
		venueName, venueLocation, venuePrice, err := rq.GetVenueNameAndPrice(payment.Reservation.VenueID)
		if err != nil {
			log.Sugar().Error("error retrieving venue name and price:", err)
			return reservation.PaymentCore{}, err
		}
		payment.Reservation.Venue.Name = venueName
		payment.Reservation.Venue.Location = venueLocation
		payment.Reservation.Venue.Price = venuePrice
	}

	return paymentToCore(payment), nil
}

// CheckAvailability implements reservation.ReservationData.
func (rq *reservationQuery) CheckAvailability(venueId string) ([]reservation.AvailabilityCore, error) {
	var result []Availability
	query := rq.db.Raw(`
	SELECT venues.venue_id,
		venues.name, 
		venues.category,
		payments.payment_id,  
		reservations.reservation_id, 
		reservations.check_in_date, 
		reservations.check_out_date
	FROM payments 
	INNER JOIN reservations ON reservations.payment_id = payments.payment_id 
	INNER JOIN venues ON venues.venue_id = reservations.venue_id
	WHERE reservations.check_in_date BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 3 DAY)
		AND payments.status IN ('success', 'pending')
		AND venues.venue_id = ?
	GROUP BY venues.venue_id, reservations.reservation_id
	`, venueId).
		Scan(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list reservations record not found")
		return nil, errors.New("list reservations record not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing list reservations query:", query.Error)
		return nil, query.Error
	} else {
		log.Sugar().Info("list reservations data found in the database")
	}

	availabilities := modelToAvailabilityCore(result)
	return availabilities, nil
}

func (rq *reservationQuery) GetReservationsByTimeSlot(venueID string, checkInDate, checkOutDate time.Time) ([]reservation.ReservationCore, error) {
	var reservations []Reservation
	query := rq.db.Where("venue_id = ? AND ((check_in_date BETWEEN ? AND ?) OR (check_out_date BETWEEN ? AND ?))",
		venueID, checkInDate, checkOutDate, checkInDate, checkOutDate).
		Find(&reservations)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list reservations record not found")
		return nil, errors.New("list reservations record not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing list reservations query:", query.Error)
		return nil, query.Error
	} else {
		log.Sugar().Info("list reservations data found in the database")
	}

	reservationCores := modelToReservationCore(reservations)
	return reservationCores, nil
}

// ReservationHistory implements reservation.ReservationData.
func (rq *reservationQuery) MyVenueCharts(userId string, keyword string, request reservation.MyReservationCore) ([]reservation.MyReservationCore, error) {
	result := []MyReservation{}
	search := "%" + keyword + "%"
	query := rq.db.Raw(`
	SELECT venues.venue_id,
		venues.name AS venue_name,
		COUNT(payments.payment_id) AS sales_volume
	FROM payments
	JOIN reservations ON payments.payment_id = reservations.payment_id
	JOIN venues ON reservations.venue_id = venues.venue_id
	WHERE reservations.user_id = ? 
		AND reservations.check_in_date BETWEEN ? AND ?
		AND payments.status LIKE ?
	GROUP BY venues.venue_id;
	`, userId, request.CheckInDate, request.CheckOutDate, search).
		Scan(&result)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list reservations record not found")
		return nil, errors.New("list reservations record not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing list reservations query:", query.Error)
		return nil, query.Error
	} else {
		log.Sugar().Info("list reservations data found in the database")
	}
	log.Sugar().Info(result)
	log.Sugar().Info(modelToMyReservationCore(result))
	return modelToMyReservationCore(result), nil
}
