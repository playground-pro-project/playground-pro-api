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
	UserID         string                    `gorm:"primaryKey;type:varchar(45)"`
	Fullname       string                    `gorm:"type:varchar(225);not null"`
	Email          string                    `gorm:"type:varchar(225);not null;unique"`
	Phone          string                    `gorm:"type:varchar(15);not null;unique"`
	Password       string                    `gorm:"type:text;not null"`
	Bio            string                    `gorm:"type:text"`
	Address        string                    `gorm:"type:text"`
	Role           string                    `gorm:"type:enum('User', 'Owner');default:'User'"`
	ProfilePicture string                    `gorm:"type:text"`
	CreatedAt      time.Time                 `gorm:"type:datetime"`
	UpdatedAt      time.Time                 `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt            `gorm:"index"`
	Venues         []venue.Venue             `gorm:"foreignKey:UserID;;foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Reservations   []reservation.Reservation `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Reviews        []review.Review           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

func UserEntityToModel(u user.UserEntity) User {
	return User{
		UserID:         u.UserID,
		Fullname:       u.Fullname,
		Email:          u.Email,
		Phone:          u.Phone,
		Password:       u.Password,
		Bio:            u.Bio,
		Address:        u.Address,
		ProfilePicture: u.ProfilePicture,
	}
}

func UserModelToEntity(u User) user.UserEntity {
	return user.UserEntity{
		UserID:         u.UserID,
		Fullname:       u.Fullname,
		Email:          u.Email,
		Phone:          u.Phone,
		Password:       u.Password,
		Bio:            u.Bio,
		Address:        u.Address,
		Role:           u.Role,
		ProfilePicture: u.ProfilePicture,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		DeletedAt:      u.DeletedAt.Time,
	}
}
