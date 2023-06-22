package data

import (
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"gorm.io/gorm"
)

var log = middlewares.Log()

type venueQuery struct {
	db *gorm.DB
}

// func New(db *gorm.DB) venue.VenueData {
// 	return &venueQuery{
// 		db: db,
// 	}
// }
