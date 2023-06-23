package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/cache"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
	"gorm.io/gorm"
)

var log = middlewares.Log()

type venueQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) venue.VenueData {
	return &venueQuery{
		db: db,
	}
}

// SearchVenue implements venue.VenueData.
func (vq *venueQuery) SearchVenue(keyword string, page pagination.Pagination) ([]venue.VenueCore, int64, int, error) {
	venues := []Venue{}
	search := "%" + keyword + "%"
	cacheKey := fmt.Sprintf("venues:%s:%d", keyword, page.Page)
	cachedVenues, err := cache.GetCached(context.Background(), cacheKey)
	if err != nil {
		return nil, 0, 0, err
	}

	if cachedVenues != nil {
		result, ok := cachedVenues.([]venue.VenueCore)
		if !ok {
			return nil, 0, 0, errors.New("unexpected type for cachedVenues")
		}
		return result, page.TotalRows, page.TotalPages, nil
	}

	query := vq.db.Table("venues").
		Select("venues.*, AVG(reviews.rating) AS average_rating, COUNT(reviews.review_id) AS total_reviews, users.fullname").
		Joins("LEFT JOIN venue_pictures ON venue_pictures.venue_id = venues.venue_id").
		Joins("LEFT JOIN reviews ON reviews.venue_id = venues.venue_id").
		Joins("LEFT JOIN users ON users.user_id = venues.owner_id").
		Where("venues.category LIKE ? AND venues.location LIKE ? AND venues.price LIKE ? AND venues.deleted_at IS NULL", search, search, search).
		Group("venues.venue_id").
		Preload("User").
		Preload("VenuePictures").
		Preload("Reviews").
		Scopes(pagination.Paginate(&venues, &page, vq.db)).
		Find(&venues)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list venues not found")
		return nil, 0, 0, errors.New("venues not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing venues query:", query.Error)
		return nil, 0, 0, query.Error
	} else {
		log.Sugar().Info("Venues data found in the database")
	}

	result := make([]venue.VenueCore, len(venues))
	for i, venue := range venues {
		result[i] = searchVenueModels(venue)
	}

	expTime := 5 * time.Second
	err = cache.SetCached(context.Background(), cacheKey, result, expTime)
	if err != nil {
		return nil, 0, 0, err
	}

	return result, page.TotalRows, page.TotalPages, nil
}
