// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	usecase "github.com/AbdulRasyid-Ans/xyz-multifinance/internal/usecase"
	mock "github.com/stretchr/testify/mock"
)

// TransactionUsecase is an autogenerated mock type for the TransactionUsecase type
type TransactionUsecase struct {
	mock.Mock
}

// CreateTransaction provides a mock function with given fields: ctx, req
func (_m *TransactionUsecase) CreateTransaction(ctx context.Context, req usecase.TransactionRequest) (usecase.GetTransactionResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for CreateTransaction")
	}

	var r0 usecase.GetTransactionResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.TransactionRequest) (usecase.GetTransactionResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.TransactionRequest) usecase.GetTransactionResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(usecase.GetTransactionResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.TransactionRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRemainingPayment provides a mock function with given fields: ctx, req
func (_m *TransactionUsecase) GetRemainingPayment(ctx context.Context, req usecase.TransactionRequest) (usecase.RemainingPaymentResponse, error) {
	ret := _m.Called(ctx, req)

	if len(ret) == 0 {
		panic("no return value specified for GetRemainingPayment")
	}

	var r0 usecase.RemainingPaymentResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, usecase.TransactionRequest) (usecase.RemainingPaymentResponse, error)); ok {
		return rf(ctx, req)
	}
	if rf, ok := ret.Get(0).(func(context.Context, usecase.TransactionRequest) usecase.RemainingPaymentResponse); ok {
		r0 = rf(ctx, req)
	} else {
		r0 = ret.Get(0).(usecase.RemainingPaymentResponse)
	}

	if rf, ok := ret.Get(1).(func(context.Context, usecase.TransactionRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewTransactionUsecase creates a new instance of TransactionUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionUsecase {
	mock := &TransactionUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
