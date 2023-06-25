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
	Username      string
	ServiceTime   time.Time
	Location      string
	Distance      uint
	Price         float64
	Longitude     float64
	Latitude      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	TotalReviews  uint
	AverageRating float64
	VenuePictures []VenuePictureCore
	Reviews       []ReviewCore
	User          UserCore
}

type ReviewCore struct {
	ReviewID  string
	UserID    string
	VenueID   string
	Review    string
	Rating    float64
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

type VenuePictureCore struct {
	VenuePictureID string
	VenueID        string
	URL            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type ReservationCore struct {
	ReservationID string
	UserID        string
	VenueID       string
	CheckInDate   time.Time
	CheckOutDate  time.Time
	Duration      uint
	Subtotal      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type VenueHandler interface {
	RegisterVenue() echo.HandlerFunc
	SearchVenue() echo.HandlerFunc
}

type VenueService interface {
	RegisterVenue(userId string, request VenueCore) (VenueCore, error)
	SearchVenue(keyword string, page pagination.Pagination) ([]VenueCore, int64, int, error)
}

type VenueData interface {
	RegisterVenue(userId string, request VenueCore) (VenueCore, error)
	SearchVenue(keyword string, page pagination.Pagination) ([]VenueCore, int64, int, error)
}
