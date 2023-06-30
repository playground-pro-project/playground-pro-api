package data

import (
	"time"

	"github.com/playground-pro-project/playground-pro-api/features/reservation"
	paymentgateway "github.com/playground-pro-project/playground-pro-api/utils/payment_gateway"
	"gorm.io/gorm"
)

type Reservation struct {
	ReservationID string `gorm:"primaryKey;type:varchar(45)"`
	UserID        string `gorm:"foreignKey:UserID;type:varchar(45)"`
	VenueID       string `gorm:"foreignKey:VenueID;type:varchar(45)"`
	PaymentID     *string
	CheckInDate   time.Time `gorm:"type:datetime"`
	CheckOutDate  time.Time `gorm:"type:datetime"`
	Duration      float64
	CreatedAt     time.Time      `gorm:"type:datetime"`
	UpdatedAt     time.Time      `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Venue         Venue          `gorm:"foreignKey:VenueID"`
}

type Payment struct {
	PaymentID     string `gorm:"primaryKey;type:varchar(45)"`
	PaymentMethod string
	PaymentType   string
	PaymentCode   string `gorm:"type:varchar(225);not null"`
	GrandTotal    string
	ServiceFee    float64
	Status        string         `gorm:"type:enum('pending','success','cancel','expire');default:'pending'"`
	CreatedAt     time.Time      `gorm:"type:datetime"`
	UpdatedAt     time.Time      `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Reservation   Reservation    `gorm:"foreignKey:PaymentID;references:PaymentID"`
}

type Venue struct {
	VenueID      string         `gorm:"primaryKey;type:varchar(45)"`
	UserID       string         `gorm:"type:varchar(45)"`
	Category     string         `gorm:"type:enum('basketball','football','futsal','badminton');default:'basketball'"`
	Name         string         `gorm:"type:varchar(225);not null;unique"`
	Description  string         `gorm:"type:text"`
	ServiceTime  string         `gorm:"type:varchar(100)"`
	Location     string         `gorm:"type:text"`
	Price        float64        `gorm:"type:double"`
	Longitude    float64        `gorm:"type:double"`
	Latitude     float64        `gorm:"type:double"`
	CreatedAt    time.Time      `gorm:"type:datetime"`
	UpdatedAt    time.Time      `gorm:"type:datetime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	Reservations []Reservation  `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

// `gorm:"type:enum('none','card','bca','bri','bni','mandiri','qris','gopay','shopeepay');default:'none'"`
// `gorm:"type:enum('cash','debit_card','bank_transfer','e-wallet');default:'cash'"`
// Struct helper for query raw in gorm
type Availability struct {
	VenueID       string
	Name          string
	Category      string
	PaymentID     string
	ReservationID string
	CheckInDate   time.Time
	CheckOutDate  time.Time
}

type MyReservation struct {
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
}

func modelToMyReservationCore(result []MyReservation) []reservation.MyReservationCore {
	myReservations := []reservation.MyReservationCore{}
	for _, r := range result {
		myReservation := reservation.MyReservationCore{
			VenueID:       r.VenueID,
			VenueName:     r.VenueName,
			Location:      r.Location,
			ReservationID: r.ReservationID,
			CheckInDate:   r.CheckInDate,
			CheckOutDate:  r.CheckOutDate,
			Duration:      r.Duration,
			Price:         r.Price,
			PaymentID:     r.PaymentID,
			GrandTotal:    r.GrandTotal,
			PaymentType:   r.PaymentType,
			PaymentCode:   r.PaymentCode,
			Status:        r.Status,
		}
		myReservations = append(myReservations, myReservation)
	}
	return myReservations
}

// Convert model to AvailabilityCore
func modelToAvailabilityCore(result []Availability) []reservation.AvailabilityCore {
	var availabilities []reservation.AvailabilityCore
	for _, r := range result {
		availability := reservation.AvailabilityCore{
			VenueID:       r.VenueID,
			Name:          r.Name,
			Category:      r.Category,
			PaymentID:     r.PaymentID,
			ReservationID: r.ReservationID,
			CheckInDate:   r.CheckInDate,
			CheckOutDate:  r.CheckOutDate,
		}
		availabilities = append(availabilities, availability)
	}
	return availabilities
}

func modelToReservationCore(models []Reservation) []reservation.ReservationCore {
	var cores []reservation.ReservationCore
	for _, m := range models {
		core := reservation.ReservationCore{
			ReservationID: m.ReservationID,
			UserID:        m.UserID,
			VenueID:       m.VenueID,
			CheckInDate:   m.CheckInDate,
			CheckOutDate:  m.CheckOutDate,
			Duration:      m.Duration,
		}
		cores = append(cores, core)
	}
	return cores
}

// Reservation-Model to reservation-core
func reservationModels(r Reservation) reservation.ReservationCore {
	return reservation.ReservationCore{
		ReservationID: r.ReservationID,
		UserID:        r.UserID,
		VenueID:       r.VenueID,
		CheckInDate:   r.CheckInDate,
		CheckOutDate:  r.CheckOutDate,
		Duration:      r.Duration,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		DeletedAt:     r.DeletedAt.Time,
	}
}

// Reservation-core to Reservation-Model
func reservationEntities(r reservation.ReservationCore) Reservation {
	return Reservation{
		ReservationID: r.ReservationID,
		UserID:        r.UserID,
		VenueID:       r.VenueID,
		CheckInDate:   r.CheckInDate,
		CheckOutDate:  r.CheckOutDate,
		Duration:      r.Duration,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		DeletedAt:     gorm.DeletedAt{Time: r.DeletedAt},
	}
}

// Payment-Model to payment-core
func paymentModels(p Payment) reservation.PaymentCore {
	return reservation.PaymentCore{
		PaymentID:     p.PaymentID,
		PaymentCode:   p.PaymentCode,
		PaymentMethod: p.PaymentMethod,
		PaymentType:   p.PaymentType,
		GrandTotal:    p.GrandTotal,
		ServiceFee:    p.ServiceFee,
		Status:        p.Status,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// Payment-core to payment-model
func paymentEntities(p reservation.PaymentCore) *Payment {
	return &Payment{
		PaymentID:     p.PaymentID,
		PaymentCode:   p.PaymentCode,
		PaymentMethod: p.PaymentMethod,
		PaymentType:   p.PaymentType,
		GrandTotal:    p.GrandTotal,
		ServiceFee:    p.ServiceFee,
		Status:        p.Status,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// Payment response Midtrans to payment-core
func PaymentCoreFromChargeResponse(res *paymentgateway.ChargeResponse) reservation.PaymentCore {
	return reservation.PaymentCore{
		PaymentID:     res.TransactionID,
		PaymentMethod: res.PaymentType,
		PaymentType:   paymentgateway.GetBankType(res),
		PaymentCode:   paymentgateway.GetPaymentCode(res),
		GrandTotal:    res.GrossAmount,
		ServiceFee:    0,
		Status:        res.TransactionStatus,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		ReservationID: res.OrderID,
		Reservation:   reservation.ReservationCore{},
	}
}

func paymentToCore(p Payment) reservation.PaymentCore {
	reservationCore := reservation.ReservationCore{
		CheckInDate:  p.Reservation.CheckInDate,
		CheckOutDate: p.Reservation.CheckOutDate,
		Duration:     p.Reservation.Duration,
	}

	paymentCore := reservation.PaymentCore{
		PaymentType: p.PaymentType,
		PaymentCode: p.PaymentCode,
		GrandTotal:  p.GrandTotal,
		ServiceFee:  p.ServiceFee,
		Status:      p.Status,
		Reservation: reservationCore,
	}

	if p.Reservation.VenueID != "" {
		paymentCore.Reservation.Venue = reservation.VenueCore{
			Name:     p.Reservation.Venue.Name,
			Location: p.Reservation.Venue.Location,
			Price:    p.Reservation.Venue.Price,
		}
	}
	return paymentCore
}
