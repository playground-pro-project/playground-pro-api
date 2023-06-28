package data

import (
	"errors"

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
		return reservation.PaymentCore{}, errors.New("no payment record has been updated")
	}

	if query.Error != nil {
		log.Error("error while updating payment status")
		return reservation.PaymentCore{}, errors.New("internal server error")
	}

	return paymentModels(*req), nil
}

// PriceVenue retrieves the price of a venue by its ID
func (rq *reservationQuery) PriceVenue(venueID string) (float64, error) {
	venue := Venue{}
	query := rq.db.Table("venues").
		Select("venues.price").
		Where("venue_id = ?", venueID).
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
