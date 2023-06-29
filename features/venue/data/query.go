package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/cache"
	"github.com/playground-pro-project/playground-pro-api/utils/helper"
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

// RegisterVenue implements venue.VenueData.
func (vq *venueQuery) RegisterVenue(userId string, request venue.VenueCore) (venue.VenueCore, error) {
	venueId := helper.GenerateVenueID()
	request.VenueID = venueId
	request.OwnerID = userId
	req := venueEntities(request)
	query := vq.db.Table("venues").Create(&req)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list venues not found")
		return venue.VenueCore{}, errors.New("venues not found")
	}
	if query.Error != nil {
		log.Error("error executing query, duplicated")
		return venue.VenueCore{}, errors.New("error executing query, duplicated")
	}
	if query.RowsAffected == 0 {
		log.Warn("no venue has been created")
		return venue.VenueCore{}, errors.New("row affected : 0")
	}

	log.Sugar().Infof("new venue has been created: %s", req.VenueID)
	return venueModels(req), nil
}

// SearchVenue implements venue.VenueData.
func (vq *venueQuery) SearchVenues(keyword string, page pagination.Pagination) ([]venue.VenueCore, int64, int, error) {
	venues := []Venue{}
	search := "%" + keyword + "%"
	expTime := 5 * time.Second
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

	err = cache.SetCached(context.Background(), cacheKey, result, expTime)
	if err != nil {
		return nil, 0, 0, err
	}

	return result, page.TotalRows, page.TotalPages, nil
}

// SelectVenueById implements venue.VenueData.
func (vq *venueQuery) SelectVenue(venueId string) (venue.VenueCore, error) {
	venues := Venue{}
	query := vq.db.Table("venues").
		Select("venues.*, AVG(reviews.rating) AS average_rating, COUNT(reviews.review_id) AS total_reviews, users.fullname").
		Joins("LEFT JOIN venue_pictures ON venue_pictures.venue_id = venues.venue_id").
		Joins("LEFT JOIN reviews ON reviews.venue_id = venues.venue_id").
		Joins("LEFT JOIN users ON users.user_id = venues.owner_id").
		Where("venues.venue_id = ?", venueId).
		Group("venues.venue_id").
		Order("venues.updated_at DESC").
		Preload("User").
		Preload("VenuePictures").
		Preload("Reviews").
		First(&venues)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list venues not found")
		return venue.VenueCore{}, errors.New("venues not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing venues query:", query.Error)
		return venue.VenueCore{}, query.Error
	} else {
		log.Sugar().Info("venues data found in the database")
	}

	return selectVenueModels(venues), nil
}

// EditVenue implements venue.VenueData.
func (vq *venueQuery) EditVenue(userId string, venueId string, request venue.VenueCore) error {
	req := venueEntities(request)
	query := vq.db.Table("venues").
		Where("owner_id = ? AND venue_id = ?", userId, venueId).
		Updates(&req)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue profile record not found")
		return errors.New("venue profile record not found")
	}

	if query.RowsAffected == 0 {
		log.Warn("no venue has been created")
		return errors.New("row affected : 0")
	}

	if query.Error != nil {
		log.Sugar().Error("error executing venues query:", query.Error)
		return errors.New("error executing venues query")
	}
	return nil
}

// UnregisterVenue implements venue.VenueData.
func (vq *venueQuery) UnregisterVenue(userId string, venueId string) error {
	query := vq.db.Table("venues").
		Where("owner_id = ? AND venue_id = ?", userId, venueId).
		Delete(&Venue{})

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue record not found")
		return errors.New("venue record not found")
	}

	if query.RowsAffected == 0 {
		log.Warn("no venue has been created")
		return errors.New("no row affected")
	}

	if query.Error != nil {
		log.Error("error while delete venue")
		return errors.New("error executing query")
	}

	return nil
}

// VenueAvailability implements venue.VenueData.
func (vq *venueQuery) VenueAvailability(venueId string) (venue.VenueCore, error) {
	venues := Venue{}
	query := vq.db.Table("venues").
		Select("venues.venue_id, venues.owner_id, venues.category, venues.name, reservations.reservation_id, reservations.check_in_date, reservations.check_out_date, users.fullname").
		Joins("JOIN reservations ON reservations.venue_id = venues.venue_id").
		Joins("JOIN users ON users.user_id = reservations.user_id").
		Where("venues.venue_id = ?", venueId).
		Where("reservations.check_in_date BETWEEN NOW() AND DATE_ADD(NOW(), INTERVAL 3 DAY)").
		Group("venues.venue_id, reservations.reservation_id").
		Preload("Reservations.User").
		First(&venues)
	if query.Error != nil {
		if query.Error == gorm.ErrRecordNotFound {
			log.Warn("venues not found")
			return venue.VenueCore{}, errors.New("venues not found")
		}
		log.Sugar().Error("error executing venues query:", query.Error)
		return venue.VenueCore{}, query.Error
	}

	result := Availability(venues)
	log.Sugar().Info(result)
	return result, nil
}

func (vq *venueQuery) InsertVenueImage(req venue.VenuePictureCore) (venue.VenuePictureCore, error) {
	venueImageID := helper.GenerateImageID()
	req.VenuePictureID = venueImageID

	model := VenuePictureCoreToModel(req)

	query := vq.db.Table("venue_pictures").Create(&model)
	if query.Error != nil {
		log.Error("error insert data, duplicated" + query.Error.Error())
		return venue.VenuePictureCore{}, errors.New("error insert data, duplicated")
	}

	rowAffect := query.RowsAffected
	if rowAffect == 0 {
		log.Warn("no venue image has been created")
		return venue.VenuePictureCore{}, errors.New("row affected : 0")
	}

	log.Sugar().Infof("new venue image has been created: %s", venueImageID)
	return VenuePictureModelToCore(model), nil
}

func (vq *venueQuery) GetAllVenueImage(venueID string) ([]venue.VenuePictureCore, error) {
	var venueImages []VenuePicture
	query := vq.db.Table("venue_pictures").Where("venue_id = ?", venueID).Find(&venueImages)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue image record not found")
		return nil, errors.New("venue image record not found")
	}

	if query.Error != nil {
		log.Error("error retrieve all images venue" + query.Error.Error())
		return nil, errors.New("error retrieve all images venue")
	}

	var venueImagesCore []venue.VenuePictureCore
	for _, img := range venueImages {
		venueImagesCore = append(venueImagesCore, VenuePictureModelToCore(img))
	}

	log.Sugar().Info(venueImagesCore)
	return venueImagesCore, nil
}

func (vq *venueQuery) DeleteVenueImage(venueImageID string) error {
	query := vq.db.Table("venue_pictures").Where("venue_picture_id = ?", venueImageID).Delete(&VenuePicture{})
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue image record not found")
		return errors.New("venue image record not found")
	}

	if query.Error != nil {
		log.Error("failed to delete image: " + query.Error.Error())
		return errors.New("failed to delete image")
	}

	return nil
}

func (vq *venueQuery) GetVenueImageByID(venueID, venueImageID string) (venue.VenuePictureCore, error) {
	var image VenuePicture
	query := vq.db.Table("venue_pictures").Where("venue_id = ? AND venue_picture_id = ?", venueID, venueImageID).First(&image)
	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("venue image record not found")
		return venue.VenuePictureCore{}, errors.New("venue image record not found")
	}

	if query.Error != nil {
		log.Error("error retrieve image venue" + query.Error.Error())
		return venue.VenuePictureCore{}, errors.New("error retrieve image venue")
	}

	imageCore := VenuePictureModelToCore(image)

	return imageCore, nil
}

// MyVenues implements venue.VenueData.
func (vq *venueQuery) MyVenues(userId string) ([]venue.VenueCore, error) {
	venues := []Venue{}
	query := vq.db.Table("venues").
		Select("venues.*, AVG(reviews.rating) AS average_rating, COUNT(reviews.review_id) AS total_reviews, users.fullname").
		Joins("LEFT JOIN venue_pictures ON venue_pictures.venue_id = venues.venue_id").
		Joins("LEFT JOIN reviews ON reviews.venue_id = venues.venue_id").
		Joins("LEFT JOIN users ON users.user_id = venues.owner_id").
		Where("venues.owner_id = ? AND venues.deleted_at IS NULL", userId).
		Group("venues.venue_id").
		Order("venues.updated_at DESC").
		Preload("User").
		Preload("VenuePictures").
		Preload("Reviews").
		Find(&venues)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		log.Error("list venues not found")
		return nil, errors.New("venues not found")
	} else if query.Error != nil {
		log.Sugar().Error("error executing venues query:", query.Error)
		return nil, query.Error
	} else {
		log.Sugar().Info("Venues data found in the database")
	}

	result := make([]venue.VenueCore, len(venues))
	for i, venue := range venues {
		result[i] = searchVenueModels(venue)
	}

	return result, nil
}
