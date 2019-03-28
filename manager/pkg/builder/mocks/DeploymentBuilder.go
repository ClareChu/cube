// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"

// DeploymentBuilder is an autogenerated mock type for the DeploymentBuilder type
type DeploymentBuilder struct {
	mock.Mock
}

// CreateApp provides a mock function with given fields: deploy
func (_m *DeploymentBuilder) CreateApp(deploy *v1alpha1.Deployment) error {
	ret := _m.Called(deploy)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1alpha1.Deployment) error); ok {
		r0 = rf(deploy)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: name, namespace, event, phase
func (_m *DeploymentBuilder) Update(name string, namespace string, event string, phase string) error {
	ret := _m.Called(name, namespace, event, phase)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string, string) error); ok {
		r0 = rf(name, namespace, event, phase)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}