// Code generated by mockery v2.53.0. DO NOT EDIT.

package mock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	models "avito-tech-merch/internal/models"
)

// TransactionRepository is an autogenerated mock type for the TransactionRepository type
type TransactionRepository struct {
	mock.Mock
}

// CreateTransaction provides a mock function with given fields: ctx, transaction
func (_m *TransactionRepository) CreateTransaction(ctx context.Context, transaction *models.Transaction) (int, error) {
	ret := _m.Called(ctx, transaction)

	if len(ret) == 0 {
		panic("no return value specified for CreateTransaction")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction) (int, error)); ok {
		return rf(ctx, transaction)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Transaction) int); ok {
		r0 = rf(ctx, transaction)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Transaction) error); ok {
		r1 = rf(ctx, transaction)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTransactionByUserID provides a mock function with given fields: ctx, userID
func (_m *TransactionRepository) GetTransactionByUserID(ctx context.Context, userID int) ([]*models.Transaction, error) {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for GetTransactionByUserID")
	}

	var r0 []*models.Transaction
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int) ([]*models.Transaction, error)); ok {
		return rf(ctx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int) []*models.Transaction); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Transaction)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
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
