// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"
import v1alpha1 "hidevops.io/mio/pkg/apis/mio/v1alpha1"

// DeploymentConfigBuilder is an autogenerated mock type for the DeploymentConfigBuilder type
type DeploymentConfigBuilder struct {
	mock.Mock
}

// CreateApp provides a mock function with given fields: deploy
func (_m *DeploymentConfigBuilder) CreateApp(deploy *v1alpha1.Deployment) error {
	ret := _m.Called(deploy)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1alpha1.Deployment) error); ok {
		r0 = rf(deploy)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
