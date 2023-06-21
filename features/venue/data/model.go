package data

import (
	"time"

	image "github.com/playground-pro-project/playground-pro-api/features/image/data"
	review "github.com/playground-pro-project/playground-pro-api/features/review/data"
	"gorm.io/gorm"
)

type Venue struct {
	VenueID       string               `gorm:"primaryKey;type:varchar(45)"`
	OwnerID       string               `gorm:"type:varchar(45)"`
	Category      string               `gorm:"type:enum('Basketball','Football','Futsal','Badminton','Swimming');default:'Basketball'"`
	Name          string               `gorm:"type:varchar(225);not null"`
	Description   string               `gorm:"type:text"`
	Location      string               `gorm:"type:text"`
	Price         float64              `gorm:"type:double"`
	Longitude     float64              `gorm:"type:double"`
	Latitude      float64              `gorm:"type:double"`
	CreatedAt     time.Time            `gorm:"type:datetime"`
	UpdatedAt     time.Time            `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt       `gorm:"index"`
	VenuePictures []image.VenuePicture `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reviews       []review.Review      `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
