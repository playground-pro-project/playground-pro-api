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
		case strings.Contains(err.Error(), "Name"):
			log.Warn("n`ame cannot be empty")
			return venue.VenueCore{}, errors.New("fullname cannot be empty")
		case strings.Contains(err.Error(), "ServiceTime"):
			log.Warn("service time cannot be empty")
			return venue.VenueCore{}, errors.New("service time cannot be empty")
		case strings.Contains(err.Error(), "Location"):
			log.Warn("location cannot be empty")
			return venue.VenueCore{}, errors.New("location cannot be empty")
		}
	}

	result, err := vs.query.RegisterVenue(userId, request)
	if err != nil {
		message := ""
		if strings.Contains(err.Error(), "error insert data, duplicated") {
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
		if strings.Contains(err.Error(), "not found") {
			log.Error("list venues record not found")
			return []venue.VenueCore{}, 0, 0, errors.New("list venues record not found")
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
func (*venueService) EditVenue(userId string, venueId string, request venue.VenueCore) error {
	panic("unimplemented")
}

// UnregisterVenue implements venue.VenueService.
func (*venueService) UnregisterVenue(userId string, venueId string) error {
	panic("unimplemented")
}

// VenueAvailability implements venue.VenueService.
func (*venueService) VenueAvailability(venueId string) ([]venue.VenueCore, error) {
	panic("unimplemented")
}
