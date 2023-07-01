package venue

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
)

type VenueCore struct {
	VenueID       string
	OwnerID       string
	Category      string `validate:"required"`
	Name          string `validate:"required"`
	Description   string
	Username      string
	ServiceTime   string `validate:"required"`
	Location      string `validate:"required"`
	Distance      float64
	Price         float64 `validate:"required"`
	Longitude     float64
	Latitude      float64
	TotalRows     int64
	TotalPages    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	TotalReviews  uint
	AverageRating float64
	VenuePictures []VenuePictureCore
	Reviews       []ReviewCore
	Reservations  []ReservationCore
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
	Username      string
	CheckInDate   time.Time
	CheckOutDate  time.Time
	Duration      uint
	Subtotal      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	User          UserCore
}

type VenueCoreRaw struct {
	VenueID       string
	OwnerID       string
	Category      string `validate:"required"`
	Name          string `validate:"required"`
	Description   string
	Username      string
	ServiceTime   string `validate:"required"`
	Location      string `validate:"required"`
	Distance      float64
	Price         float64 `validate:"required"`
	Longitude     float64
	Latitude      float64
	TotalRows     int64
	TotalPages    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
	TotalReviews  uint
	AverageRating float64
	VenuePicture  string
	VenuePictures []VenuePictureCore
	Reviews       []ReviewCore
	Reservations  []ReservationCore
	User          UserCore
}

type VenueHandler interface {
	// RegisterVenue() echo.HandlerFunc
	SearchVenues() echo.HandlerFunc
	SelectVenue() echo.HandlerFunc
	EditVenue() echo.HandlerFunc
	UnregisterVenue() echo.HandlerFunc
	VenueAvailability() echo.HandlerFunc
	// CreateVenueImage() echo.HandlerFunc
	GetAllVenueImage() echo.HandlerFunc
	DeleteVenueImage() echo.HandlerFunc
	MyVenues() echo.HandlerFunc
	CreateVenue() echo.HandlerFunc
}

type VenueService interface {
	RegisterVenue(userId string, request VenueCore) (VenueCore, error)
	SearchVenues(keyword string, latitude float64, longitude float64, page pagination.Pagination) ([]VenueCoreRaw, int64, int, error)
	SelectVenue(venueId string) (VenueCore, error)
	EditVenue(userId string, venueId string, request VenueCore) error
	UnregisterVenue(userId string, venueId string) error
	VenueAvailability(venueId string) (VenueCore, error)
	// CreateVenueImage(req VenuePictureCore) (VenuePictureCore, error)
	GetAllVenueImage(venueID string) ([]VenuePictureCore, error)
	DeleteVenueImage(venueImageID string) error
	GetVenueImageByID(venueID, venueImageID string) (VenuePictureCore, error)
	MyVenues(userId string) ([]VenueCore, error)
	CreateVenue(userID string, venueReq VenueCore, venueImageReq VenuePictureCore) (VenueCore, error)
}

type VenueData interface {
	RegisterVenue(userId string, request VenueCore) (VenueCore, error)
	SearchVenues(keyword string, latitude float64, longitude float64, page pagination.Pagination) ([]VenueCoreRaw, int64, int, error)
	SelectVenue(venueId string) (VenueCore, error)
	EditVenue(userId string, venueId string, request VenueCore) error
	UnregisterVenue(userId string, venueId string) error
	VenueAvailability(venueId string) (VenueCore, error)
	// InsertVenueImage(req VenuePictureCore) (VenuePictureCore, error)
	GetAllVenueImage(venueID string) ([]VenuePictureCore, error)
	DeleteVenueImage(venueImageID string) error
	GetVenueImageByID(venueID, venueImageID string) (VenuePictureCore, error)
	MyVenues(userId string) ([]VenueCore, error)
	InsertVenue(userID string, venueReq VenueCore, venueImageReq VenuePictureCore) (VenueCore, error)
}
