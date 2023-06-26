package data

import (
	"time"

	image "github.com/playground-pro-project/playground-pro-api/features/image/data"
	reservation "github.com/playground-pro-project/playground-pro-api/features/reservation/data"
	review "github.com/playground-pro-project/playground-pro-api/features/review/data"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"gorm.io/gorm"
)

type Venue struct {
	VenueID       string                    `gorm:"primaryKey;type:varchar(45)"`
	OwnerID       string                    `gorm:"type:varchar(45)"`
	Category      string                    `gorm:"type:enum('Basketball','Football','Futsal','Badminton','Swimming');default:'Basketball'"`
	Name          string                    `gorm:"type:varchar(225);not null"`
	Description   string                    `gorm:"type:text"`
	ServiceTime   string                    `gorm:"type:varchar(100)"`
	Location      string                    `gorm:"type:text"`
	Price         float64                   `gorm:"type:double"`
	Longitude     float64                   `gorm:"type:double"`
	Latitude      float64                   `gorm:"type:double"`
	CreatedAt     time.Time                 `gorm:"type:datetime"`
	UpdatedAt     time.Time                 `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt            `gorm:"index"`
	User          User                      `gorm:"references:OwnerID;foreignKey:UserID"`
	VenuePictures []image.VenuePicture      `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reservations  []reservation.Reservation `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reviews       []review.Review           `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type User struct {
	UserID         string                    `gorm:"primaryKey;type:varchar(45)"`
	Fullname       string                    `gorm:"type:varchar(225);not null"`
	Email          string                    `gorm:"type:varchar(225);not null;unique"`
	Phone          string                    `gorm:"type:varchar(15);not null;unique"`
	Password       string                    `gorm:"type:text;not null"`
	Bio            string                    `gorm:"type:text"`
	Address        string                    `gorm:"type:text"`
	Role           string                    `gorm:"type:enum('user', 'owner', 'admin');default:'user'"`
	ProfilePicture string                    `gorm:"type:text"`
	CreatedAt      time.Time                 `gorm:"type:datetime"`
	UpdatedAt      time.Time                 `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt            `gorm:"index"`
	Venues         []Venue                   `gorm:"foreignKey:UserID;;foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reservations   []reservation.Reservation `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reviews        []review.Review           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

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

type VenuePicture struct {
	VenuePictureID string         `gorm:"primaryKey;type:varchar(45)"`
	VenueID        string         `gorm:"type:varchar(45)"`
	URL            string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"type:datetime"`
	UpdatedAt      time.Time      `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func searchVenueModels(v Venue) venue.VenueCore {
	var reviews []venue.ReviewCore
	var totalRating float64
	var averageRating float64 = 0.0
	var picture string
	for _, r := range v.Reviews {
		tmp := venue.ReviewCore{
			Review: r.Review,
			Rating: r.Rating,
		}
		reviews = append(reviews, tmp)
		totalRating += r.Rating
	}

	if len(v.Reviews) > 0 {
		averageRating = totalRating / float64(len(v.Reviews))
	}

	if len(v.VenuePictures) > 0 {
		picture = v.VenuePictures[0].URL
	}

	result := venue.VenueCore{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Username:      v.User.Fullname,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Price:         v.Price,
		Longitude:     v.Longitude,
		Latitude:      v.Latitude,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		DeletedAt:     v.DeletedAt.Time,
		TotalReviews:  uint(len(v.Reviews)),
		AverageRating: averageRating,
		VenuePictures: []venue.VenuePictureCore{
			{
				URL: picture,
			},
		},
	}

	return result
}

func selectVenueModels(v Venue) venue.VenueCore {
	var reviews []venue.ReviewCore
	var totalRating float64
	var averageRating float64 = 0.0
	for _, r := range v.Reviews {
		tmp := venue.ReviewCore{
			Review: r.Review,
			Rating: r.Rating,
		}
		reviews = append(reviews, tmp)
		totalRating += r.Rating
	}

	if len(v.Reviews) > 0 {
		averageRating = totalRating / float64(len(v.Reviews))
	}

	pictures := make([]venue.VenuePictureCore, len(v.VenuePictures))
	for i, p := range v.VenuePictures {
		pictures[i] = venue.VenuePictureCore{
			URL: p.URL,
		}
	}

	result := venue.VenueCore{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		Username:      v.User.Fullname,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Price:         v.Price,
		Longitude:     v.Longitude,
		Latitude:      v.Latitude,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		DeletedAt:     v.DeletedAt.Time,
		TotalReviews:  uint(len(v.Reviews)),
		AverageRating: averageRating,
		VenuePictures: pictures,
		Reviews:       reviews,
	}

	return result
}

// Venue-Model to venue-core
func venueModels(v Venue) venue.VenueCore {
	return venue.VenueCore{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Price:         v.Price,
		Longitude:     v.Longitude,
		Latitude:      v.Latitude,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		DeletedAt:     v.DeletedAt.Time,
		VenuePictures: []venue.VenuePictureCore{},
		Reviews:       []venue.ReviewCore{},
	}
}

// Venue-core to venue-model
func venueEntities(v venue.VenueCore) Venue {
	return Venue{
		VenueID:       v.VenueID,
		OwnerID:       v.OwnerID,
		Category:      v.Category,
		Name:          v.Name,
		Description:   v.Description,
		ServiceTime:   v.ServiceTime,
		Location:      v.Location,
		Price:         v.Price,
		Longitude:     v.Longitude,
		Latitude:      v.Latitude,
		CreatedAt:     v.CreatedAt,
		UpdatedAt:     v.UpdatedAt,
		DeletedAt:     gorm.DeletedAt{Time: v.DeletedAt},
		VenuePictures: []image.VenuePicture{},
		Reviews:       []review.Review{},
	}
}
