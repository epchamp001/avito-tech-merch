// Code generated by mockery v2.53.0. DO NOT EDIT.

package mock

import (
	db "avito-tech-merch/internal/storage/db"
	context "context"

	mock "github.com/stretchr/testify/mock"

	pgx "github.com/jackc/pgx/v5"
)

// TxManager is an autogenerated mock type for the TxManager type
type TxManager struct {
	mock.Mock
}

// GetExecutor provides a mock function with given fields: ctx
func (_m *TxManager) GetExecutor(ctx context.Context) db.Executor {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetExecutor")
	}

	var r0 db.Executor
	if rf, ok := ret.Get(0).(func(context.Context) db.Executor); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(db.Executor)
		}
	}

	return r0
}

// WithTx provides a mock function with given fields: ctx, isoLevel, accessMode, fn
func (_m *TxManager) WithTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode, fn func(context.Context) error) error {
	ret := _m.Called(ctx, isoLevel, accessMode, fn)

	if len(ret) == 0 {
		panic("no return value specified for WithTx")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pgx.TxIsoLevel, pgx.TxAccessMode, func(context.Context) error) error); ok {
		r0 = rf(ctx, isoLevel, accessMode, fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTxManager creates a new instance of TxManager. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTxManager(t interface {
	mock.TestingT
	Cleanup(func())
}) *TxManager {
	mock := &TxManager{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
