package user

import (
	"time"
)

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
	Reservations   []ReservationEntity
	Reviews        []ReviewEntity
}

type VenueEntity struct {
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
	VenuePictures []VenuePictureEntity
	Reviews       []ReviewEntity
}

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

type ReservationEntity struct {
	ReservationID string
	UserID        string
	VenueVenueID  string
	PaymentID     string
	CheckInDate   time.Time
	CheckOutDate  time.Time
	Duration      uint
	Subtotal      float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     time.Time
}

type VenuePictureEntity struct {
	VenuePictureID string
	VenueID        string
	URL            string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      time.Time
}

type UserData interface {
	Create(user UserEntity) (string, error)
	GetByID(userID string) (UserEntity, error)
	GetAll() ([]UserEntity, error)
	UpdateByID(userID string, updatedUser UserEntity) error
	DeleteByID(userID string) error
	Login(email, password string) (UserEntity, string, error)
}

type UserService interface {
	CreateUser(user UserEntity) (string, error)
	GetUserByID(userID string) (UserEntity, error)
	GetAllUsers() ([]UserEntity, error)
	UpdateUserByID(userID string, updatedUser UserEntity) error
	DeleteUserByID(userID string) error
	Login(email, password string) (UserEntity, string, error)
}
