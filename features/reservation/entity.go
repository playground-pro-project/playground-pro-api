package user

import "time"

type ReservationCore struct {
<<<<<<< HEAD
	ReservationID string    
	UserID        string    
	VenueVenueID  string    
	PaymentID     string    
	CheckInDate   time.Time 
	CheckOutDate  time.Time 
=======
	ReservationID string
	UserID        string
	VenueVenueID  string
	PaymentID     string
	CheckInDate   time.Time
	CheckOutDate  time.Time
>>>>>>> e733440 (Update entity to core)
	Duration      uint
	Subtotal      float64
	CreatedAt     time.Time      
	UpdatedAt     time.Time      
	DeletedAt     time.Time
}
