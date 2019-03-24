// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import command "hidevops.io/cube/console/pkg/command"
import mock "github.com/stretchr/testify/mock"
import v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"

// ServiceConfigAggregate is an autogenerated mock type for the ServiceConfigAggregate type
type ServiceConfigAggregate struct {
	mock.Mock
}

// Create provides a mock function with given fields: name, pipelineName, namespace, sourceType, version, profile
func (_m *ServiceConfigAggregate) Create(name string, pipelineName string, namespace string, sourceType string, version string, profile string) (*v1alpha1.ServiceConfig, error) {
	ret := _m.Called(name, pipelineName, namespace, sourceType, version, profile)

	var r0 *v1alpha1.ServiceConfig
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, string) *v1alpha1.ServiceConfig); ok {
		r0 = rf(name, pipelineName, namespace, sourceType, version, profile)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.ServiceConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, string, string, string, string) error); ok {
		r1 = rf(name, pipelineName, namespace, sourceType, version, profile)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteService provides a mock function with given fields: name, namespace
func (_m *ServiceConfigAggregate) DeleteService(name string, namespace string) error {
	ret := _m.Called(name, namespace)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(name, namespace)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Template provides a mock function with given fields: cmd
func (_m *ServiceConfigAggregate) Template(cmd *command.ServiceConfig) (*v1alpha1.ServiceConfig, error) {
	ret := _m.Called(cmd)

	var r0 *v1alpha1.ServiceConfig
	if rf, ok := ret.Get(0).(func(*command.ServiceConfig) *v1alpha1.ServiceConfig); ok {
		r0 = rf(cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.ServiceConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*command.ServiceConfig) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
