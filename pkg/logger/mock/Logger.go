// Code generated by mockery v2.53.0. DO NOT EDIT.

package mock

import (
	mock "github.com/stretchr/testify/mock"
	zapcore "go.uber.org/zap/zapcore"
)

// Logger is an autogenerated mock type for the Logger type
type Logger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: msg, fields
func (_m *Logger) Debug(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Debugw provides a mock function with given fields: msg, keysAndValues
func (_m *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// Error provides a mock function with given fields: msg, fields
func (_m *Logger) Error(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Errorw provides a mock function with given fields: msg, keysAndValues
func (_m *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// Fatal provides a mock function with given fields: msg, fields
func (_m *Logger) Fatal(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Fatalw provides a mock function with given fields: msg, keysAndValues
func (_m *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// Info provides a mock function with given fields: msg, fields
func (_m *Logger) Info(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Infow provides a mock function with given fields: msg, keysAndValues
func (_m *Logger) Infow(msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// Sync provides a mock function with no fields
func (_m *Logger) Sync() {
	_m.Called()
}

// Warn provides a mock function with given fields: msg, fields
func (_m *Logger) Warn(msg string, fields ...zapcore.Field) {
	_va := make([]interface{}, len(fields))
	for _i := range fields {
		_va[_i] = fields[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, _va...)
	_m.Called(_ca...)
}

// Warnw provides a mock function with given fields: msg, keysAndValues
func (_m *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, keysAndValues...)
	_m.Called(_ca...)
}

// NewLogger creates a new instance of Logger. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewLogger(t interface {
	mock.TestingT
	Cleanup(func())
}) *Logger {
	mock := &Logger{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
