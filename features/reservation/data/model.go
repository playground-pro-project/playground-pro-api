package data

import (
	"time"

	payment "github.com/playground-pro-project/playground-pro-api/features/payment/data"
	"gorm.io/gorm"
)

type Reservation struct {
	ReservationID string         `gorm:"primaryKey;type:varchar(45)"`
	UserID        string         `gorm:"primaryKey;type:varchar(45)"`
	VenueID  	  string         `gorm:"primaryKey;type:varchar(45)"`
	CheckInDate   time.Time      `gorm:"type:datetime"`
	CheckOutDate  time.Time      `gorm:"type:datetime"`
	Duration      uint
	Subtotal      float64       
	CreatedAt     time.Time      `gorm:"type:datetime"`
	UpdatedAt     time.Time      `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Payment payment.Payment `gorm:"foreignKey:ReservationID"`
}
