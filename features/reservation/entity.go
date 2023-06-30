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
	ReservationID string
	Reservation   ReservationCore
	Venue         VenueCore
}

type VenueCore struct {
	VenueID      string
	OwnerID      string
	Category     string
	Name         string
	Description  string
	Username     string
	ServiceTime  string
	Location     string
	Distance     uint
	Price        float64
	Longitude    float64
	Latitude     float64
	Reservations []ReservationCore
}

type AvailabilityCore struct {
	VenueID       string
	Name          string
	Category      string
	PaymentID     string
	ReservationID string
	CheckInDate   time.Time
	CheckOutDate  time.Time
}

type MyReservationCore struct {
	VenueID       string
	VenueName     string
	Location      string
	ReservationID string
	CheckInDate   time.Time
	CheckOutDate  time.Time
	Duration      float64
	Price         float64
	PaymentID     string
	GrandTotal    float64
	PaymentType   string
	PaymentCode   string
	Status        string
	SalesVolume   uint
}

type ReservationHandler interface {
	MakeReservation() echo.HandlerFunc
	ReservationStatus() echo.HandlerFunc
	MyReservation() echo.HandlerFunc
	DetailTransaction() echo.HandlerFunc
	CheckAvailability() echo.HandlerFunc
	MyVenueCharts() echo.HandlerFunc
}

type ReservationService interface {
	MakeReservation(userId string, r ReservationCore, p PaymentCore) (ReservationCore, PaymentCore, error)
	ReservationStatus(request PaymentCore) (PaymentCore, error)
	MyReservation(userId string) ([]MyReservationCore, error)
	DetailTransaction(userId string, paymentId string) (PaymentCore, error)
	CheckAvailability(venueId string) ([]AvailabilityCore, error)
	MyVenueCharts(userId string, keyword string, request MyReservationCore) ([]MyReservationCore, error)
}

type ReservationData interface {
	MakeReservation(userId string, r ReservationCore, p PaymentCore) (ReservationCore, PaymentCore, error)
	ReservationStatus(request PaymentCore) (PaymentCore, error)
	PriceVenue(venueID string) (float64, error)
	ReservationCheckOutDate(reservation_id string) (time.Time, error)
	MyReservation(userId string) ([]MyReservationCore, error)
	DetailTransaction(userId string, paymentId string) (PaymentCore, error)
	CheckAvailability(venueId string) ([]AvailabilityCore, error)
	GetReservationsByTimeSlot(venueID string, checkInDate, checkOutDate time.Time) ([]ReservationCore, error)
	MyVenueCharts(userId string, keyword string, request MyReservationCore) ([]MyReservationCore, error)
}
