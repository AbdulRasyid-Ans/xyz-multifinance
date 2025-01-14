// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	repository "github.com/AbdulRasyid-Ans/xyz-multifinance/internal/repository"
	mock "github.com/stretchr/testify/mock"

	sql "database/sql"
)

// TransactionRepository is an autogenerated mock type for the TransactionRepository type
type TransactionRepository struct {
	mock.Mock
}

// BeginTx provides a mock function with given fields: ctx
func (_m *TransactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for BeginTx")
	}

	var r0 *sql.Tx
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*sql.Tx, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *sql.Tx); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Tx)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CommitTx provides a mock function with given fields: ctx, tx
func (_m *TransactionRepository) CommitTx(ctx context.Context, tx *sql.Tx) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for CommitTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CreateTransaction provides a mock function with given fields: ctx, tx, transaction
func (_m *TransactionRepository) CreateTransaction(ctx context.Context, tx *sql.Tx, transaction repository.Transaction) (int64, error) {
	ret := _m.Called(ctx, tx, transaction)

	if len(ret) == 0 {
		panic("no return value specified for CreateTransaction")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, repository.Transaction) (int64, error)); ok {
		return rf(ctx, tx, transaction)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx, repository.Transaction) int64); ok {
		r0 = rf(ctx, tx, transaction)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *sql.Tx, repository.Transaction) error); ok {
		r1 = rf(ctx, tx, transaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByID provides a mock function with given fields: ctx, id
func (_m *TransactionRepository) GetTransactionByID(ctx context.Context, id int64) (repository.Transaction, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionByID")
	}

	var r0 repository.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (repository.Transaction, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) repository.Transaction); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(repository.Transaction)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionsByConsumerID provides a mock function with given fields: ctx, consumerID
func (_m *TransactionRepository) GetTransactionsByConsumerID(ctx context.Context, consumerID int64) ([]repository.Transaction, error) {
	ret := _m.Called(ctx, consumerID)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionsByConsumerID")
	}

	var r0 []repository.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]repository.Transaction, error)); ok {
		return rf(ctx, consumerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []repository.Transaction); ok {
		r0 = rf(ctx, consumerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, consumerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionsByLoanID provides a mock function with given fields: ctx, loanID
func (_m *TransactionRepository) GetTransactionsByLoanID(ctx context.Context, loanID int64) ([]repository.Transaction, error) {
	ret := _m.Called(ctx, loanID)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionsByLoanID")
	}

	var r0 []repository.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) ([]repository.Transaction, error)); ok {
		return rf(ctx, loanID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) []repository.Transaction); ok {
		r0 = rf(ctx, loanID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]repository.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, loanID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RollbackTx provides a mock function with given fields: ctx, tx
func (_m *TransactionRepository) RollbackTx(ctx context.Context, tx *sql.Tx) error {
	ret := _m.Called(ctx, tx)

	if len(ret) == 0 {
		panic("no return value specified for RollbackTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *sql.Tx) error); ok {
		r0 = rf(ctx, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTransactionRepository creates a new instance of TransactionRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTransactionRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *TransactionRepository {
	mock := &TransactionRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
