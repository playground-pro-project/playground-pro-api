package data

import (
	"time"

	"github.com/playground-pro-project/playground-pro-api/features/review"
	"gorm.io/gorm"
)

type Review struct {
	ReviewID  string         `gorm:"primaryKey;type:varchar(45)"`
	UserID    string         `gorm:"type:varchar(45)"`
	VenueID   string         `gorm:"type:varchar(45)"`
	Review    string         `gorm:"type:text"`
	Rating    float64        `gorm:"type:double"`
	CreatedAt time.Time      `gorm:"type:datetime"`
	UpdatedAt time.Time      `gorm:"type:datetime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"references:UserID"`
	Venue     Venue          `gorm:"references:VenueID"`
}

type User struct {
	UserID         string   `gorm:"primaryKey;type:varchar(45)"`
	Fullname       string   `gorm:"type:varchar(225);not null"`
	Email          string   `gorm:"type:varchar(225);not null;unique"`
	Phone          string   `gorm:"type:varchar(15);not null;unique"`
	Password       string   `gorm:"type:text;not null"`
	Bio            string   `gorm:"type:text"`
	Address        string   `gorm:"type:text"`
	Role           string   `gorm:"type:enum('user', 'owner', 'admin');default:'user'"`
	ProfilePicture string   `gorm:"type:text"`
	Reviews        []Review `gorm:"foreignKey:UserID"`
	Venues         []Venue  `gorm:"foreignKey:VenueID"`
}

type Venue struct {
	VenueID     string         `gorm:"primaryKey;type:varchar(45)"`
	OwnerID     string         `gorm:"type:varchar(45)"`
	Category    string         `gorm:"type:enum('Basketball','Football','Futsal','Badminton','Swimming');default:'Basketball'"`
	Name        string         `gorm:"type:varchar(225);not null"`
	Description string         `gorm:"type:text"`
	Location    string         `gorm:"type:text"`
	Price       float64        `gorm:"type:double"`
	Longitude   float64        `gorm:"type:double"`
	Latitude    float64        `gorm:"type:double"`
	CreatedAt   time.Time      `gorm:"type:datetime"`
	UpdatedAt   time.Time      `gorm:"type:datetime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Reviews     []Review       `gorm:"foreignKey:ReviewID"`
}

func ReviewCoreToModel(r review.ReviewCore) Review {
	return Review{
		UserID:  r.UserID,
		VenueID: r.VenueID,
		Review:  r.Review,
		Rating:  r.Rating,
	}
}

func ReviewModelToCore(r Review) review.ReviewCore {
	return review.ReviewCore{
		ReviewID:  r.ReviewID,
		UserID:    r.UserID,
		VenueID:   r.VenueID,
		Review:    r.Review,
		Rating:    r.Rating,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
		DeletedAt: r.DeletedAt.Time,
		User:      UserModelToCore(r.User),
	}
}

func UserModelToCore(u User) review.UserCore {
	return review.UserCore{
		UserID:         u.UserID,
		Fullname:       u.Fullname,
		Email:          u.Email,
		Phone:          u.Phone,
		Password:       u.Password,
		Bio:            u.Bio,
		Address:        u.Address,
		Role:           u.Role,
		ProfilePicture: u.ProfilePicture,
	}
}

func VenueModelToCore(v Venue) review.VenueCore {
	return review.VenueCore{
		VenueID:     v.VenueID,
		OwnerID:     v.OwnerID,
		Category:    v.Category,
		Name:        v.Name,
		Description: v.Description,
		Location:    v.Location,
		Price:       v.Price,
		Longitude:   v.Longitude,
		Latitude:    v.Latitude,
	}
}
