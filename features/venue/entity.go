package venue

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
)

type VenueCore struct {
	VenueID       string 
	OwnerID       string 
	Category      string 
	Name          string 
	Description   string 
	Location      string 
	Price         float64
	Longitude     float64
	Latitude      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	VenuePictures []VenuePictureCore
	Reviews       []ReviewCore
}

type VenuePictureCore struct {
	VenuePictureID string
	VenueID        string
	URL            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
