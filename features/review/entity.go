package review

import "time"

type ReviewEntity struct {
	ReviewID  string
	UserID    string
	VenueID   string
	Review    string
	Rating    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	User      UserEntity
	Venue     VenueEntity
}

type UserEntity struct {
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
	Venues         []VenueEntity
}

type VenueEntity struct {
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
	Reviews     []ReviewEntity
}
