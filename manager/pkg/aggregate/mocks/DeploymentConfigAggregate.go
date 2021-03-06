// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import command "hidevops.io/cube/manager/pkg/command"
import mock "github.com/stretchr/testify/mock"
import v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"

// DeploymentConfigAggregate is an autogenerated mock type for the DeploymentConfigAggregate type
type DeploymentConfigAggregate struct {
	mock.Mock
}

// Create provides a mock function with given fields: param, buildVersion
func (_m *DeploymentConfigAggregate) Create(param *command.PipelineReqParams, buildVersion string) (*v1alpha1.DeploymentConfig, error) {
	ret := _m.Called(param, buildVersion)

	var r0 *v1alpha1.DeploymentConfig
	if rf, ok := ret.Get(0).(func(*command.PipelineReqParams, string) *v1alpha1.DeploymentConfig); ok {
		r0 = rf(param, buildVersion)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.DeploymentConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*command.PipelineReqParams, string) error); ok {
		r1 = rf(param, buildVersion)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// InitDeployConfig provides a mock function with given fields: deploy, template, param
func (_m *DeploymentConfigAggregate) InitDeployConfig(deploy *v1alpha1.DeploymentConfig, template *v1alpha1.DeploymentConfig, param *command.PipelineReqParams) {
	_m.Called(deploy, template, param)
}

// Template provides a mock function with given fields: cmd
func (_m *DeploymentConfigAggregate) Template(cmd *command.DeploymentConfig) (*v1alpha1.DeploymentConfig, error) {
	ret := _m.Called(cmd)

	var r0 *v1alpha1.DeploymentConfig
	if rf, ok := ret.Get(0).(func(*command.DeploymentConfig) *v1alpha1.DeploymentConfig); ok {
		r0 = rf(cmd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.DeploymentConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*command.DeploymentConfig) error); ok {
		r1 = rf(cmd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
