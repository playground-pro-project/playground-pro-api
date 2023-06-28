package service

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/playground-pro-project/playground-pro-api/app/middlewares"
	"github.com/playground-pro-project/playground-pro-api/features/venue"
	"github.com/playground-pro-project/playground-pro-api/utils/pagination"
)

var log = middlewares.Log()

type venueService struct {
	query    venue.VenueData
	validate *validator.Validate
}

func New(vd venue.VenueData) venue.VenueService {
	return &venueService{
		query:    vd,
		validate: validator.New(),
	}
}

// RegisterVenue implements venue.VenueService.
func (vs *venueService) RegisterVenue(userId string, request venue.VenueCore) (venue.VenueCore, error) {
	err := vs.validate.Struct(request)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "category"):
			log.Warn("category cannot be empty")
			return venue.VenueCore{}, errors.New("category cannot be empty")
		case strings.Contains(err.Error(), "name"):
			log.Warn("name cannot be empty")
			return venue.VenueCore{}, errors.New("name cannot be empty")
		case strings.Contains(err.Error(), "servicetime"):
			log.Warn("service time cannot be empty")
			return venue.VenueCore{}, errors.New("service time cannot be empty")
		case strings.Contains(err.Error(), "location"):
			log.Warn("location cannot be empty")
			return venue.VenueCore{}, errors.New("location cannot be empty")
		case strings.Contains(err.Error(), "price"):
			log.Warn("price cannot be empty")
			return venue.VenueCore{}, errors.New("price cannot be empty")
		}
	}

	result, err := vs.query.RegisterVenue(userId, request)
	if err != nil {
		message := ""
		if strings.Contains(err.Error(), "duplicated") {
			log.Error("error insert data, duplicated")
			message = "error insert data, duplicated"
		} else {
			log.Error("internal server error")
			message = "internal server error"
		}
		return venue.VenueCore{}, errors.New(message)
	}

	return result, nil
}

// SearchVenue implements venue.VenueService.
func (vs *venueService) SearchVenues(keyword string, page pagination.Pagination) ([]venue.VenueCore, int64, int, error) {
	if page.Sort != "" {
		ps := strings.Replace(page.Sort, "_", " ", 1)
		page.Sort = ps
	}

	venues, rows, pages, err := vs.query.SearchVenues(keyword, page)
	if err != nil {
		if strings.Contains(err.Error(), "venues not found") {
			log.Error("list venues record not found")
			return []venue.VenueCore{}, 0, 0, errors.New("venues not found")
		} else {
			log.Error("internal server error")
			return []venue.VenueCore{}, 0, 0, errors.New("internal server error")
		}
	}

	return venues, rows, pages, err
}

// SelectVenue implements venue.VenueService.
func (vs *venueService) SelectVenue(venueId string) (venue.VenueCore, error) {
	result, err := vs.query.SelectVenue(venueId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Error("not found, error while retrieving venue")
			return venue.VenueCore{}, errors.New("not found, error while retrieving venue")
		} else {
			log.Error("internal server error")
			return venue.VenueCore{}, errors.New("internal server error")
		}
	}
	return result, nil
}

// EditVenue implements venue.VenueService.
func (vs *venueService) EditVenue(userId string, venueId string, request venue.VenueCore) error {
	err := vs.query.EditVenue(userId, venueId, request)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Error("venue profile record not found")
			return errors.New("venue profile record not found")
		} else if strings.Contains(err.Error(), "duplicate data entry") {
			log.Error("failed to update venue, duplicate data entry")
			return errors.New("failed to update venue, duplicate data entry")
		} else {
			log.Error("internal server error")
			return errors.New("internal server error")
		}
	}

	return nil
}

// UnregisterVenue implements venue.VenueService.
func (vs *venueService) UnregisterVenue(userId string, venueId string) error {
	err := vs.query.UnregisterVenue(userId, venueId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Error("venue record not found")
			return errors.New("venue record not found")
		} else {
			log.Error("internal server error")
			return errors.New("internal server error")
		}
	}

	return nil
}

// VenueAvailability implements venue.VenueService.
func (vs *venueService) VenueAvailability(venueId string) (venue.VenueCore, error) {
	venues, err := vs.query.VenueAvailability(venueId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.Warn("list venues record not found")
			return venue.VenueCore{}, errors.New("list venues record not found")
		} else {
			log.Error("internal server error")
			return venue.VenueCore{}, errors.New("internal server error")
		}
	}

	return venues, err
}

func (vs *venueService) CreateVenueImage(req venue.VenuePictureCore) (venue.VenuePictureCore, error) {
	switch {
	case req.VenueID == "":
		log.Error("error, venue ID is required")
		return venue.VenuePictureCore{}, errors.New("error, venue ID is required")
	case req.URL == "":
		log.Error("error, venue URL image is required")
		return venue.VenuePictureCore{}, errors.New("error, venue URL image is required")
	}

	vn, err := vs.query.InsertVenueImage(req)
	if err != nil {
		log.Error(err.Error())
		return venue.VenuePictureCore{}, err
	}

	return vn, nil
}

func (vs *venueService) GetAllVenueImage(venueID string) ([]venue.VenuePictureCore, error) {
	venueImages, err := vs.query.GetAllVenueImage(venueID)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return venueImages, nil
}

func (vs *venueService) DeleteVenueImage(venueImageID string) error {
	err := vs.query.DeleteVenueImage(venueImageID)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (vs *venueService) GetByVenueImageID(venueImageID string) (venue.VenuePictureCore, error) {
	venueImage, err := vs.query.GetByVenueImageID(venueImageID)
	if err != nil {
		log.Error(err.Error())
		return venue.VenuePictureCore{}, err
	}

	return venueImage, nil
}
