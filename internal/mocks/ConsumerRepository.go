// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	repository "github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	mock "github.com/stretchr/testify/mock"
)

// ConsumerRepository is an autogenerated mock type for the ConsumerRepository type
type ConsumerRepository struct {
	mock.Mock
}

// CreateConsumer provides a mock function with given fields: ctx, consumer
func (_m *ConsumerRepository) CreateConsumer(ctx context.Context, consumer repository.Consumer) (int64, error) {
	ret := _m.Called(ctx, consumer)

	if len(ret) == 0 {
		panic("no return value specified for CreateConsumer")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Consumer) (int64, error)); ok {
		return rf(ctx, consumer)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.Consumer) int64); ok {
		r0 = rf(ctx, consumer)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.Consumer) error); ok {
		r1 = rf(ctx, consumer)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteConsumer provides a mock function with given fields: ctx, id
func (_m *ConsumerRepository) DeleteConsumer(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteConsumer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FetchConsumer provides a mock function with given fields: ctx, req
func (_m *ConsumerRepository) FetchConsumer(ctx context.Context, req repository.FetchConsumerRequest) ([]repository.Consumer, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for FetchConsumer")
	}

	var r0 []repository.Consumer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.FetchConsumerRequest) ([]repository.Consumer, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repository.FetchConsumerRequest) []repository.Consumer); ok {
		r0 = rf(ctx, req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Consumer)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repository.FetchConsumerRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConsumerByID provides a mock function with given fields: ctx, id
func (_m *ConsumerRepository) GetConsumerByID(ctx context.Context, id int64) (repository.Consumer, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetConsumerByID")
	}

	var r0 repository.Consumer
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (repository.Consumer, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) repository.Consumer); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(repository.Consumer)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateConsumer provides a mock function with given fields: ctx, consumer
func (_m *ConsumerRepository) UpdateConsumer(ctx context.Context, consumer repository.Consumer) error {
	ret := _m.Called(ctx, consumer)

	if len(ret) == 0 {
		panic("no return value specified for UpdateConsumer")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, repository.Consumer) error); ok {
		r0 = rf(ctx, consumer)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewConsumerRepository creates a new instance of ConsumerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewConsumerRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ConsumerRepository {
	mock := &ConsumerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
