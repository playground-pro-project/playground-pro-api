package reservation

import (
	"time"

	"github.com/labstack/echo/v4"
)

type ReservationCore struct {
	ReservationID string
	UserID        string
	VenueID       string    `validate:"required"`
	CheckInDate   time.Time `validate:"required"`
	CheckOutDate  time.Time `validate:"required"`
	Duration      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	Venue         VenueCore
}

type PaymentCore struct {
	PaymentID     string
	PaymentMethod string
	PaymentType   string
	PaymentCode   string
	GrandTotal    string
	ServiceFee    float64
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Reservation   ReservationCore
}

type VenueCore struct {
	VenueID      string
	Name         string
	Location     string
	Price        float64
	Reservations []ReservationCore
}

type ReservationHandler interface {
	MakeReservation() echo.HandlerFunc
	ReservationStatus() echo.HandlerFunc
	ReservationHistory() echo.HandlerFunc
	DetailTransaction() echo.HandlerFunc
}

type ReservationService interface {
	MakeReservation(userId string, r ReservationCore, p PaymentCore) (ReservationCore, PaymentCore, error)
	ReservationStatus(request PaymentCore) (PaymentCore, error)
	ReservationHistory(userId string) ([]PaymentCore, error)
	DetailTransaction(userId string, paymentId string) (PaymentCore, error)
}

type ReservationData interface {
	MakeReservation(userId string, r ReservationCore, p PaymentCore) (ReservationCore, PaymentCore, error)
	ReservationStatus(request PaymentCore) (PaymentCore, error)
	PriceVenue(venueID string) (float64, error)
	ReservationCheckOutDate(reservation_id string) (time.Time, error)
	ReservationHistory(userId string) ([]PaymentCore, error)
	DetailTransaction(userId string, paymentId string) (PaymentCore, error)
}
