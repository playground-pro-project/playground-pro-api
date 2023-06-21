package data

import (
	"time"

	image "github.com/playground-pro-project/playground-pro-api/features/image/data"
	review "github.com/playground-pro-project/playground-pro-api/features/review/data"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"gorm.io/gorm"
)

type Venue struct {
	VenueID       string `gorm:"primaryKey;type:varchar(45)"`
	OwnerID       string `gorm:"type:varchar(45)"`
	Category      string `gorm:"type:enum('Basketball','Football','Futsal','Badminton','Swimming');default:'Basketball'"`
	Name          string `gorm:"type:varchar(225);not null"`
	Description   string `gorm:"type:text"`
	Location      string `gorm:"type:text"`
	Price         float64
	Longitude     float64
	Latitude      float64
	CreatedAt     time.Time            `gorm:"type:datetime"`
	UpdatedAt     time.Time            `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt       `gorm:"index"`
	VenuePictures []image.VenuePicture `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reviews       []review.Review      `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}



// Venue-Model to venue-core
func venueModels(v Venue) venue.VenueCore{
	return venue.VenueCore{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Location:      v.Location,
		Price:         v.Price,
		Longitude:     v.Longitude,
		Latitude:      v.Latitude,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		DeletedAt:     v.DeletedAt.Time,
		VenuePictures: []venue.VenuePictureCore{},
		Reviews:       []venue.ReviewCore{},
	}
}

// Venue-core to venue-model
func venueEntities(v venue.VenueCore) Venue {
	return Venue{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Location:      v.Location,
		Price:         v.Price,
		Longitude:     v.Longitude,
		Latitude:      v.Latitude,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		DeletedAt:     gorm.DeletedAt{Time: v.DeletedAt},
		VenuePictures: []image.VenuePicture{},
		Reviews:       []review.Review{},
	}
}
