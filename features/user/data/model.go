package data

import (
	"time"

	reservation "github.com/playground-pro-project/playground-pro-api/features/reservation/data"
	review "github.com/playground-pro-project/playground-pro-api/features/review/data"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"gorm.io/gorm"
)

type User struct {
	UserID         string                    `gorm:"primaryKey;type:varchar(45)"`
	Fullname       string                    `gorm:"type:varchar(225);not null"`
	Email          string                    `gorm:"type:varchar(225);not null;unique"`
	Phone          string                    `gorm:"type:varchar(15);not null;unique"`
	Password       string                    `gorm:"type:text;not null"`
	Bio            string                    `gorm:"type:text"`
	Address        string                    `gorm:"type:text"`
	Role           string                    `gorm:"type:enum('user', 'owner', 'admin');default:'user'"`
	AccountStatus  string                    `gorm:"type:enum('verified', 'unverified');default:'unverified'"`
	ProfilePicture string                    `gorm:"type:varchar(255);default:'https://cdn.pixabay.com/photo/2015/10/05/22/37/blank-profile-picture-973460_1280.png'"`
	OwnerFile      string                    `gorm:"type:text"`
	CreatedAt      time.Time                 `gorm:"type:datetime"`
	UpdatedAt      time.Time                 `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt            `gorm:"index"`
	Venues         []Venue                   `gorm:"foreignKey:UserID;foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reservations   []reservation.Reservation `gorm:"foreignKey:UserID"`
	Reviews        []review.Review           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Venue struct {
	VenueID       string          `gorm:"primaryKey;type:varchar(45)"`
	OwnerID       string          `gorm:"type:varchar(45)"`
	Category      string          `gorm:"type:enum('basketball','football','futsal','badminton');default:'basketball'"`
	Name          string          `gorm:"type:varchar(225);not null;unique"`
	Description   string          `gorm:"type:text"`
	ServiceTime   string          `gorm:"type:varchar(100)"`
	Location      string          `gorm:"type:text"`
	Price         float64         `gorm:"type:double"`
	Longitude     float64         `gorm:"type:double"`
	Latitude      float64         `gorm:"type:double"`
	CreatedAt     time.Time       `gorm:"type:datetime"`
	UpdatedAt     time.Time       `gorm:"type:datetime"`
	DeletedAt     gorm.DeletedAt  `gorm:"index"`
	User          User            `gorm:"references:OwnerID;foreignKey:UserID"`
	VenuePictures []VenuePicture  `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reviews       []review.Review `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type VenuePicture struct {
	VenuePictureID string         `gorm:"primaryKey;type:varchar(45)"`
	VenueID        string         `gorm:"type:varchar(45)"`
	URL            string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"type:datetime"`
	UpdatedAt      time.Time      `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func VenueModelToCore(v Venue) venue.VenueCore {
	return venue.VenueCore{
		VenueID:     v.VenueID,
		OwnerID:     v.OwnerID,
		Category:    v.Category,
		Name:        v.Name,
		Description: v.Description,
		Location:    v.Location,
		Price:       v.Price,
		Longitude:   v.Longitude,
		Latitude:    v.Latitude,
		CreatedAt:   v.CreatedAt,
		UpdatedAt:   v.UpdatedAt,
		DeletedAt:   v.DeletedAt.Time,
	}
}

func UserCoreToModel(u user.UserCore) User {
	return User{
		UserID:         u.UserID,
		Fullname:       u.Fullname,
		Email:          u.Email,
		Phone:          u.Phone,
		Password:       u.Password,
		Bio:            u.Bio,
		Address:        u.Address,
		Role:           u.Role,
		AccountStatus:  u.AccountStatus,
		ProfilePicture: u.ProfilePicture,
		OwnerFile:      u.OwnerFile,
	}
}

func UserModelToCore(u User) user.UserCore {
	return user.UserCore{
		UserID:         u.UserID,
		Fullname:       u.Fullname,
		Email:          u.Email,
		Phone:          u.Phone,
		Password:       u.Password,
		Bio:            u.Bio,
		Address:        u.Address,
		Role:           u.Role,
		AccountStatus:  u.AccountStatus,
		ProfilePicture: u.ProfilePicture,
		OwnerFile:      u.OwnerFile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt.Time,
	}
}
