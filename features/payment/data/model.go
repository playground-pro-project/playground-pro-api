package data

import (
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	PaymentID   string `gorm:"primaryKey;type:varchar(45)"`
	VANumber    string `gorm:"type:varchar(225);not null"`
	BankAccount string `gorm:"type:enum('cash','BCA','BRI','BNI','Mandiri');default:'cash'"`
	GrandTotal  float64
	ServiceFee 	float64
	Status      string         `gorm:"type:enum('Pending','Success','Cancelled','Expired');default:'Pending'"`
	CreatedAt   time.Time      `gorm:"type:datetime"`
	UpdatedAt   time.Time      `gorm:"type:datetime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
