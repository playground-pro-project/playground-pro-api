package data

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
	ReservationID string         `gorm:"primaryKey;type:varchar(45)"`
	UserID        string         `gorm:"type:varchar(45)"`
	VenueVenueID  string         `gorm:"primaryKey;type:varchar(45)"`
	PaymentID     string         `gorm:"primaryKey;type:varchar(45)"`
	CheckInDate   time.Time      `gorm:"type:datetime"`
	CheckOutDate  time.Time      `gorm:"type:datetime"`
	Duration      uint           `gorm:"type:tinyint"`
	Subtotal      float64        `gorm:"type:double"`
	CreatedAt     time.Time      `gorm:"type:datetime"`
	UpdatedAt     time.Time      `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
