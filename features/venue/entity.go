package venue

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
)

type VenueCore struct {
<<<<<<< HEAD
	VenueID       string 
	OwnerID       string 
	Category      string 
	Name          string 
	Description   string 
	Location      string 
=======
	VenueID       string
	OwnerID       string
	Category      string
	Name          string
	Description   string
	Location      string
>>>>>>> e733440 (Update entity to core)
	Price         float64
	Longitude     float64
	Latitude      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	VenuePictures []VenuePictureCore
	Reviews       []ReviewCore
}

<<<<<<< HEAD
=======
type ReviewCore struct {
	ReviewID  string
	UserID    string
	VenueID   string
	Review    string
	Rating    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	User      UserCore
	Venue     VenueCore
}

type UserCore struct {
	UserID         string
	Fullname       string
	Email          string
	Phone          string
	Password       string
	Bio            string
	Address        string
	Role           string
	ProfilePicture string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	Venues         []VenueCore
}

>>>>>>> e733440 (Update entity to core)
type VenuePictureCore struct {
	VenuePictureID string
	VenueID        string
	URL            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}
