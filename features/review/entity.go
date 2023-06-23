package review

import (
	"time"
)

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
}

type VenueCore struct {
	VenueID     string
	OwnerID     string
	Category    string
	Name        string
	Description string
	Location    string
	Price       float64
	Longitude   float64
	Latitude    float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
}

type ReviewData interface {
	Create(venueID string, userID string, review ReviewCore) (string, error)
	GetAllByVenueID(venueID string) ([]ReviewCore, error)
	DeleteByID(reviewID string) error
}

type ReviewService interface {
	CreateReview(venueID string, userID string, review ReviewCore) (string, error)
	GetAllByVenueID(venueID string) ([]ReviewCore, error)
	DeleteByID(reviewID string) error
}
