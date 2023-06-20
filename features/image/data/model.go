package data

import (
	"time"

	"gorm.io/gorm"
)

type VenuePicture struct {
	VenuePictureID string         `gorm:"primaryKey;type:varchar(45)"`
	VenueID        string         `gorm:"type:varchar(45)"`
	URL            string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"type:datetime"`
	UpdatedAt      time.Time      `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}
