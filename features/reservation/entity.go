package reservation

import (
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ReservationCore struct {
	ReservationID string
	UserID        string
	VenueID       string    `validate:"required"`
	CheckInDate   time.Time `validate:"required"`
	CheckOutDate  time.Time `validate:"required"`
	Duration      uint
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type PaymentCore struct {
	PaymentID     string
	VANumber      string
	PaymentMethod string
	PaymentType   string
	GrandTotal    string
	ServiceFee    float64
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt
	Reservation   ReservationCore
}

type ReservationHandler interface {
	MakeReservation() echo.HandlerFunc
}

type ReservationService interface {
	MakeReservation(userId string, r ReservationCore, p PaymentCore) (ReservationCore, PaymentCore, error)
}

type ReservationData interface {
	MakeReservation(userId string, r ReservationCore, p PaymentCore) (ReservationCore, PaymentCore, error)
}
