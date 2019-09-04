// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import command "hidevops.io/cube/manager/pkg/command"
import mock "github.com/stretchr/testify/mock"

import v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"

// AppService is an autogenerated mock type for the AppService type
type AppService struct {
	mock.Mock
}

// Create provides a mock function with given fields: cmd
func (_m *AppService) Create(cmd *command.PipelineStart) (*v1alpha1.App, error) {
	ret := _m.Called(cmd)

	var r0 *v1alpha1.App
	if rf, ok := ret.Get(0).(func(*command.PipelineStart) *v1alpha1.App); ok {
		r0 = rf(cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.App)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*command.PipelineStart) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: name, namespace
func (_m *AppService) Delete(name string, namespace string) error {
	ret := _m.Called(name, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: name, namespace
func (_m *AppService) Get(name string, namespace string) (*v1alpha1.App, error) {
	ret := _m.Called(name, namespace)

	var r0 *v1alpha1.App
	if rf, ok := ret.Get(0).(func(string, string) *v1alpha1.App); ok {
		r0 = rf(name, namespace)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.App)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(name, namespace)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Init provides a mock function with given fields: cmd
func (_m *AppService) Init(cmd *command.PipelineStart) (bool, error) {
	ret := _m.Called(cmd)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*command.PipelineStart) bool); ok {
		r0 = rf(cmd)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*command.PipelineStart) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: name, namespace, cmd
func (_m *AppService) Update(name string, namespace string, cmd *command.PipelineStart) (*v1alpha1.App, error) {
	ret := _m.Called(name, namespace, cmd)

	var r0 *v1alpha1.App
	if rf, ok := ret.Get(0).(func(string, string, *command.PipelineStart) *v1alpha1.App); ok {
		r0 = rf(name, namespace, cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.App)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *command.PipelineStart) error); ok {
		r1 = rf(name, namespace, cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
