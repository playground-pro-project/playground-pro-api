package data

import (
	"time"

	reservation "github.com/playground-pro-project/playground-pro-api/features/reservation/data"
	review "github.com/playground-pro-project/playground-pro-api/features/review/data"
	"github.com/playground-pro-project/playground-pro-api/features/user"
	venue "github.com/playground-pro-project/playground-pro-api/features/venue/data"
	"gorm.io/gorm"
)

type User struct {
	UserID         string `gorm:"primaryKey;type:varchar(45)"`
	Fullname       string `gorm:"type:varchar(225);not null"`
	Email          string `gorm:"type:varchar(225);not null;unique"`
	Phone          string `gorm:"type:varchar(15);not null;unique"`
	Password       string `gorm:"type:text;not null"`
	Bio            string `gorm:"type:text"`
	Address        string `gorm:"type:text"`
	Role           string `gorm:"type:enum('user', 'owner', 'admin');default:'user'"`
	ProfilePicture string `gorm:"type:text"`
	OtpEnabled     bool   `gorm:"default:false;"`
	OtpVerified    bool   `gorm:"default:false;"`
	OtpSecret      string
	OtpAuthURL     string
	OwnerFile      string                    `gorm:"type:text"`
	CreatedAt      time.Time                 `gorm:"type:datetime"`
	UpdatedAt      time.Time                 `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt            `gorm:"index"`
	Venues         []venue.Venue             `gorm:"foreignKey:UserID;foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reservations   []reservation.Reservation `gorm:"foreignKey:UserID"`
	Reviews        []review.Review           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
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
		ProfilePicture: u.ProfilePicture,
		OtpEnabled:     u.OtpEnabled,
		OtpVerified:    u.OtpVerified,
		OtpSecret:      u.OtpSecret,
		OtpAuthURL:     u.OtpAuthURL,
		OwnerFile:      u.OwnerFile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt.Time,
	}
}
