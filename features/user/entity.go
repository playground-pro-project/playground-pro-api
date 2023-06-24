package user

import (
	"time"
)

type UserCore struct {
	UserID         string
	Fullname       string `validate:"required"`
	Email          string `validate:"required,email"`
	Phone          string `validate:"required"`
	Password       string `validate:"required"`
	Bio            string
	Address        string
	Role           string
	ProfilePicture string
	OtpEnabled     bool
	OtpVerified    bool
	OtpSecret      string
	OtpAuthURL     string
	OwnerFile      string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
	Venues         []VenueCore
	Reservations   []ReservationCore
	Reviews        []ReviewCore
}

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

type VenuePictureCore struct {
	VenuePictureID string
	VenueID        string
	URL            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type UserService interface {
	Register(request UserCore) (UserCore, error)
	Login(request UserCore) (UserCore, string, error)
	GenerateOTP(request UserCore) (UserCore, error)
	DeleteByID(userID string) error
	GetByID(userID string) (UserCore, error)
	UpdateByID(userID string, updatedUser UserCore) error
}

type UserData interface {
	Register(request UserCore) (UserCore, error)
	Login(request UserCore) (UserCore, string, error)
	GenerateOTP(request UserCore) (UserCore, error)
	DeleteByID(userID string) error
	GetByID(userID string) (UserCore, error)
	UpdateByID(userID string, updatedUser UserCore) error
}
