package data

import (
	"time"

	"gorm.io/gorm"
)

type Review struct {
	ReviewID  string `gorm:"primaryKey;type:varchar(45)"`
	UserID    string `gorm:"type:varchar(45)"`
	VenueID   string `gorm:"type:varchar(45)"`
	Review    string `gorm:"type:text"`
	Rating    float64
	CreatedAt time.Time      `gorm:"type:datetime"`
	UpdatedAt time.Time      `gorm:"type:datetime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      User           `gorm:"references:UserID"`
	Venue     Venue          `gorm:"references:VenueID"`
}

type User struct {
	UserID         string         `gorm:"primaryKey;type:varchar(45)"`
	Fullname       string         `gorm:"type:varchar(225);not null"`
	Email          string         `gorm:"type:varchar(225);not null;unique"`
	Phone          string         `gorm:"type:varchar(15);not null;unique"`
	Password       string         `gorm:"type:text;not null"`
	Bio            string         `gorm:"type:text"`
	Address        string         `gorm:"type:text"`
	Role           string         `gorm:"type:enum('User', 'Owner');default:'User'"`
	ProfilePicture string         `gorm:"type:text"`
	CreatedAt      time.Time      `gorm:"type:datetime"`
	UpdatedAt      time.Time      `gorm:"type:datetime"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	Venues         []Venue        `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

type Venue struct {
	VenueID     string `gorm:"primaryKey;type:varchar(45)"`
	UserID      string `gorm:"type:varchar(45)"`
	Category    string `gorm:"type:enum('Basketball','Football','Futsal','Badminton','Swimming');default:'Basketball'"`
	Name        string `gorm:"type:varchar(225);not null"`
	Description string `gorm:"type:text"`
	Location    string `gorm:"type:text"`
	Price       float64
	Longitude   float64
	Latitude    float64
	CreatedAt   time.Time      `gorm:"type:datetime"`
	UpdatedAt   time.Time      `gorm:"type:datetime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Reviews     []Review       `gorm:"foreignKey:VenueID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}
