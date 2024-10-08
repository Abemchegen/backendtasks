// Code generated by mockery v2.44.1. DO NOT EDIT.

package mocks

import (
	domain "task8/domain"

	mock "github.com/stretchr/testify/mock"
)

// TaskRepositoryInterface is an autogenerated mock type for the TaskRepositoryInterface type
type TaskRepositoryInterface struct {
	mock.Mock
}

// CreateTask provides a mock function with given fields: newtask, userid
func (_m *TaskRepositoryInterface) CreateTask(newtask *domain.Task, userid string) error {
	ret := _m.Called(newtask, userid)

	if len(ret) == 0 {
		panic("no return value specified for CreateTask")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*domain.Task, string) error); ok {
		r0 = rf(newtask, userid)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetTask provides a mock function with given fields: id
func (_m *TaskRepositoryInterface) GetTask(id string) (*domain.Task, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetTask")
	}

	var r0 *domain.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*domain.Task, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *domain.Task); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTasks provides a mock function with given fields: userid
func (_m *TaskRepositoryInterface) GetTasks(userid string) (*[]domain.Task, error) {
	ret := _m.Called(userid)

	if len(ret) == 0 {
		panic("no return value specified for GetTasks")
	}

	var r0 *[]domain.Task
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*[]domain.Task, error)); ok {
		return rf(userid)
	}
	if rf, ok := ret.Get(0).(func(string) *[]domain.Task); ok {
		r0 = rf(userid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]domain.Task)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RemoveTask provides a mock function with given fields: id
func (_m *TaskRepositoryInterface) RemoveTask(id string) error {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for RemoveTask")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateTask provides a mock function with given fields: id, updatedtask
func (_m *TaskRepositoryInterface) UpdateTask(id string, updatedtask *domain.Task) error {
	ret := _m.Called(id, updatedtask)

	if len(ret) == 0 {
		panic("no return value specified for UpdateTask")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *domain.Task) error); ok {
		r0 = rf(id, updatedtask)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewTaskRepositoryInterface creates a new instance of TaskRepositoryInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewTaskRepositoryInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *TaskRepositoryInterface {
	mock := &TaskRepositoryInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
