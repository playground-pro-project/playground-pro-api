// Code generated by mockery v2.30.1. DO NOT EDIT.

package mocks

import (
	venue "github.com/playground-pro-project/playground-pro-api/features/venue"
	pagination "github.com/playground-pro-project/playground-pro-api/utils/pagination"
	mock "github.com/stretchr/testify/mock"
)

// VenueService is an autogenerated mock type for the VenueService type
type VenueService struct {
	mock.Mock
}

// RegisterVenue provides a mock function with given fields: userId, request
func (_m *VenueService) RegisterVenue(userId string, request venue.VenueCore) (venue.VenueCore, error) {
	ret := _m.Called(userId, request)

	var r0 venue.VenueCore
	var r1 error
	if rf, ok := ret.Get(0).(func(string, venue.VenueCore) (venue.VenueCore, error)); ok {
		return rf(userId, request)
	}
	if rf, ok := ret.Get(0).(func(string, venue.VenueCore) venue.VenueCore); ok {
		r0 = rf(userId, request)
	} else {
		r0 = ret.Get(0).(venue.VenueCore)
	}

	if rf, ok := ret.Get(1).(func(string, venue.VenueCore) error); ok {
		r1 = rf(userId, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SearchVenue provides a mock function with given fields: keyword, page
func (_m *VenueService) SearchVenue(keyword string, page pagination.Pagination) ([]venue.VenueCore, int64, int, error) {
	ret := _m.Called(keyword, page)

	var r0 []venue.VenueCore
	var r1 int64
	var r2 int
	var r3 error
	if rf, ok := ret.Get(0).(func(string, pagination.Pagination) ([]venue.VenueCore, int64, int, error)); ok {
		return rf(keyword, page)
	}
	if rf, ok := ret.Get(0).(func(string, pagination.Pagination) []venue.VenueCore); ok {
		r0 = rf(keyword, page)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]venue.VenueCore)
		}
	}

	if rf, ok := ret.Get(1).(func(string, pagination.Pagination) int64); ok {
		r1 = rf(keyword, page)
	} else {
		r1 = ret.Get(1).(int64)
	}

	if rf, ok := ret.Get(2).(func(string, pagination.Pagination) int); ok {
		r2 = rf(keyword, page)
	} else {
		r2 = ret.Get(2).(int)
	}

	if rf, ok := ret.Get(3).(func(string, pagination.Pagination) error); ok {
		r3 = rf(keyword, page)
	} else {
		r3 = ret.Error(3)
	}

	return r0, r1, r2, r3
}

// NewVenueService creates a new instance of VenueService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewVenueService(t interface {
	mock.TestingT
	Cleanup(func())
}) *VenueService {
	mock := &VenueService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
