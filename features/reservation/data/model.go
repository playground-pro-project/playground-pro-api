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
	ReservationID string
	Reservation   Reservation `gorm:"foreignKey:PaymentID;references:PaymentID"`
}

type Venue struct {
	VenueID      string         `gorm:"primaryKey;type:varchar(45)"`
	OwnerID      string         `gorm:"type:varchar(45)"`
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

func reservationToCore(r Reservation) reservation.ReservationCore {
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
		// Payment: reservation.Payment{
		// 	PaymentID:     *r.PaymentID,
		// 	PaymentMethod: r.Payment.PaymentMethod,
		// 	PaymentType:   r.Payment.PaymentType,
		// 	PaymentCode:   r.Payment.PaymentCode,
		// 	GrandTotal:    r.Payment.GrandTotal,
		// 	ServiceFee:    r.Payment.ServiceFee,
		// 	Status:        r.Payment.Status,
		// 	CreatedAt:     r.Payment.CreatedAt,
		// 	UpdatedAt:     r.Payment.UpdatedAt,
		// },
	}
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
		Reservation:   reservation.ReservationCore{},
	}
}

// `gorm:"type:enum('none','card','bca','bri','bni','mandiri','qris','gopay','shopeepay');default:'none'"`
// `gorm:"type:enum('cash','debit_card','bank_transfer','e-wallet');default:'cash'"`
