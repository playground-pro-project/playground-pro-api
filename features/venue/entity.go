package venue

import (
	"time"
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
	Reservations   []ReservationCore
	Reviews        []ReviewCore
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
