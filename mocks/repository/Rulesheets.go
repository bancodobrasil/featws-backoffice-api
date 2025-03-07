// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	models "github.com/bancodobrasil/featws-api/models"

	repository "github.com/bancodobrasil/featws-api/repository"
)

// Rulesheets is an autogenerated mock type for the Rulesheets type
type Rulesheets struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx, entity
func (_m *Rulesheets) Count(ctx context.Context, entity interface{}) (int64, error) {
	ret := _m.Called(ctx, entity)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) int64); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(ctx, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountInTransaction provides a mock function with given fields: ctx, db, entity
func (_m *Rulesheets) CountInTransaction(ctx context.Context, db *gorm.DB, entity interface{}) (int64, error) {
	ret := _m.Called(ctx, db, entity)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, interface{}) int64); ok {
		r0 = rf(ctx, db, entity)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, interface{}) error); ok {
		r1 = rf(ctx, db, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Create provides a mock function with given fields: ctx, entity
func (_m *Rulesheets) Create(ctx context.Context, entity *models.Rulesheet) error {
	ret := _m.Called(ctx, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Rulesheet) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateInTransaction provides a mock function with given fields: ctx, db, entity
func (_m *Rulesheets) CreateInTransaction(ctx context.Context, db *gorm.DB, entity *models.Rulesheet) error {
	ret := _m.Called(ctx, db, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *models.Rulesheet) error); ok {
		r0 = rf(ctx, db, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Rulesheets) Delete(ctx context.Context, id string) (bool, error) {
	ret := _m.Called(ctx, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteInTransaction provides a mock function with given fields: ctx, db, id
func (_m *Rulesheets) DeleteInTransaction(ctx context.Context, db *gorm.DB, id string) (bool, error) {
	ret := _m.Called(ctx, db, id)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) bool); ok {
		r0 = rf(ctx, db, id)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, string) error); ok {
		r1 = rf(ctx, db, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Find provides a mock function with given fields: ctx, entity, options
func (_m *Rulesheets) Find(ctx context.Context, entity interface{}, options *repository.FindOptions) ([]*models.Rulesheet, error) {
	ret := _m.Called(ctx, entity, options)

	var r0 []*models.Rulesheet
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, *repository.FindOptions) []*models.Rulesheet); ok {
		r0 = rf(ctx, entity, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Rulesheet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, *repository.FindOptions) error); ok {
		r1 = rf(ctx, entity, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindInTransaction provides a mock function with given fields: ctx, db, entity, options
func (_m *Rulesheets) FindInTransaction(ctx context.Context, db *gorm.DB, entity interface{}, options *repository.FindOptions) ([]*models.Rulesheet, error) {
	ret := _m.Called(ctx, db, entity, options)

	var r0 []*models.Rulesheet
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, interface{}, *repository.FindOptions) []*models.Rulesheet); ok {
		r0 = rf(ctx, db, entity, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Rulesheet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, interface{}, *repository.FindOptions) error); ok {
		r1 = rf(ctx, db, entity, options)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: ctx, id
func (_m *Rulesheets) Get(ctx context.Context, id string) (*models.Rulesheet, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Rulesheet
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Rulesheet); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Rulesheet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDB provides a mock function with given fields:
func (_m *Rulesheets) GetDB() *gorm.DB {
	ret := _m.Called()

	var r0 *gorm.DB
	if rf, ok := ret.Get(0).(func() *gorm.DB); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*gorm.DB)
		}
	}

	return r0
}

// GetInTransaction provides a mock function with given fields: ctx, db, id
func (_m *Rulesheets) GetInTransaction(ctx context.Context, db *gorm.DB, id string) (*models.Rulesheet, error) {
	ret := _m.Called(ctx, db, id)

	var r0 *models.Rulesheet
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) *models.Rulesheet); ok {
		r0 = rf(ctx, db, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Rulesheet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, string) error); ok {
		r1 = rf(ctx, db, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, entity
func (_m *Rulesheets) Update(ctx context.Context, entity models.Rulesheet) (*models.Rulesheet, error) {
	ret := _m.Called(ctx, entity)

	var r0 *models.Rulesheet
	if rf, ok := ret.Get(0).(func(context.Context, models.Rulesheet) *models.Rulesheet); ok {
		r0 = rf(ctx, entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Rulesheet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Rulesheet) error); ok {
		r1 = rf(ctx, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateInTransaction provides a mock function with given fields: ctx, db, entity
func (_m *Rulesheets) UpdateInTransaction(ctx context.Context, db *gorm.DB, entity models.Rulesheet) (*models.Rulesheet, error) {
	ret := _m.Called(ctx, db, entity)

	var r0 *models.Rulesheet
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, models.Rulesheet) *models.Rulesheet); ok {
		r0 = rf(ctx, db, entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Rulesheet)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, models.Rulesheet) error); ok {
		r1 = rf(ctx, db, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRulesheets interface {
	mock.TestingT
	Cleanup(func())
}

// NewRulesheets creates a new instance of Rulesheets. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRulesheets(t mockConstructorTestingTNewRulesheets) *Rulesheets {
	mock := &Rulesheets{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
