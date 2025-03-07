// Code generated by mockery. DO NOT EDIT.

package mocks

import (
	context "context"

	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"

	repository "github.com/bancodobrasil/featws-api/repository"
)

// Repository is an autogenerated mock type for the Repository type
type Repository[T interface{}] struct {
	mock.Mock
}

// Count provides a mock function with given fields: ctx, entity
func (_m *Repository[T]) Count(ctx context.Context, entity interface{}) (int64, error) {
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
func (_m *Repository[T]) CountInTransaction(ctx context.Context, db *gorm.DB, entity interface{}) (int64, error) {
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
func (_m *Repository[T]) Create(ctx context.Context, entity *T) error {
	ret := _m.Called(ctx, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *T) error); ok {
		r0 = rf(ctx, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateInTransaction provides a mock function with given fields: ctx, db, entity
func (_m *Repository[T]) CreateInTransaction(ctx context.Context, db *gorm.DB, entity *T) error {
	ret := _m.Called(ctx, db, entity)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, *T) error); ok {
		r0 = rf(ctx, db, entity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: ctx, id
func (_m *Repository[T]) Delete(ctx context.Context, id string) (bool, error) {
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
func (_m *Repository[T]) DeleteInTransaction(ctx context.Context, db *gorm.DB, id string) (bool, error) {
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
func (_m *Repository[T]) Find(ctx context.Context, entity interface{}, options *repository.FindOptions) ([]*T, error) {
	ret := _m.Called(ctx, entity, options)

	var r0 []*T
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, *repository.FindOptions) []*T); ok {
		r0 = rf(ctx, entity, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*T)
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
func (_m *Repository[T]) FindInTransaction(ctx context.Context, db *gorm.DB, entity interface{}, options *repository.FindOptions) ([]*T, error) {
	ret := _m.Called(ctx, db, entity, options)

	var r0 []*T
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, interface{}, *repository.FindOptions) []*T); ok {
		r0 = rf(ctx, db, entity, options)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*T)
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
func (_m *Repository[T]) Get(ctx context.Context, id string) (*T, error) {
	ret := _m.Called(ctx, id)

	var r0 *T
	if rf, ok := ret.Get(0).(func(context.Context, string) *T); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
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
func (_m *Repository[T]) GetDB() *gorm.DB {
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
func (_m *Repository[T]) GetInTransaction(ctx context.Context, db *gorm.DB, id string) (*T, error) {
	ret := _m.Called(ctx, db, id)

	var r0 *T
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, string) *T); ok {
		r0 = rf(ctx, db, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
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
func (_m *Repository[T]) Update(ctx context.Context, entity T) (*T, error) {
	ret := _m.Called(ctx, entity)

	var r0 *T
	if rf, ok := ret.Get(0).(func(context.Context, T) *T); ok {
		r0 = rf(ctx, entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, T) error); ok {
		r1 = rf(ctx, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateInTransaction provides a mock function with given fields: ctx, db, entity
func (_m *Repository[T]) UpdateInTransaction(ctx context.Context, db *gorm.DB, entity T) (*T, error) {
	ret := _m.Called(ctx, db, entity)

	var r0 *T
	if rf, ok := ret.Get(0).(func(context.Context, *gorm.DB, T) *T); ok {
		r0 = rf(ctx, db, entity)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*T)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *gorm.DB, T) error); ok {
		r1 = rf(ctx, db, entity)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository[T interface{}](t mockConstructorTestingTNewRepository) *Repository[T] {
	mock := &Repository[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
