// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// PasswordService is an autogenerated mock type for the PasswordService type
type PasswordService struct {
	mock.Mock
}

// Compare provides a mock function with given fields: password1, password2
func (_m *PasswordService) Compare(password1 string, password2 string) error {
	ret := _m.Called(password1, password2)

	if len(ret) == 0 {
		panic("no return value specified for Compare")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(password1, password2)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Hash provides a mock function with given fields: Password
func (_m *PasswordService) Hash(Password string) (string, error) {
	ret := _m.Called(Password)

	if len(ret) == 0 {
		panic("no return value specified for Hash")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (string, error)); ok {
		return rf(Password)
	}
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(Password)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(Password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewPasswordService creates a new instance of PasswordService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewPasswordService(t interface {
	mock.TestingT
	Cleanup(func())
}) *PasswordService {
	mock := &PasswordService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
