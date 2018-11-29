// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import bufio "bufio"
import corev1 "k8s.io/api/core/v1"
import mock "github.com/stretchr/testify/mock"
import service "hidevops.io/mio/console/pkg/service"
import v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
import v1alpha1 "hidevops.io/mio/pkg/apis/mio/v1alpha1"

// KubeClient is an autogenerated mock type for the KubeClient type
type KubeClient struct {
	mock.Mock
}

// GetBuildConfigLastVersion provides a mock function with given fields: namespace, name
func (_m *KubeClient) GetBuildConfigLastVersion(namespace string, name string) (int, error) {
	ret := _m.Called(namespace, name)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, string) int); ok {
		r0 = rf(namespace, name)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLogs provides a mock function with given fields: namespace, name, tail
func (_m *KubeClient) GetLogs(namespace string, name string, tail int64) (*bufio.Reader, error) {
	ret := _m.Called(namespace, name, tail)

	var r0 *bufio.Reader
	if rf, ok := ret.Get(0).(func(string, string, int64) *bufio.Reader); ok {
		r0 = rf(namespace, name, tail)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*bufio.Reader)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int64) error); ok {
		r1 = rf(namespace, name, tail)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPipelineApp provides a mock function with given fields: namespace, name, new, profile, version, podMessage
func (_m *KubeClient) GetPipelineApp(namespace string, name string, new string, profile string, version string, podMessage chan service.PodMessage) (string, string, error) {
	ret := _m.Called(namespace, name, new, profile, version, podMessage)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string, string, string, string, chan service.PodMessage) string); ok {
		r0 = rf(namespace, name, new, profile, version, podMessage)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(string, string, string, string, string, chan service.PodMessage) string); ok {
		r1 = rf(namespace, name, new, profile, version, podMessage)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, string, string, string, chan service.PodMessage) error); ok {
		r2 = rf(namespace, name, new, profile, version, podMessage)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetPipelineConfig provides a mock function with given fields: namespace, name
func (_m *KubeClient) GetPipelineConfig(namespace string, name string) (*v1alpha1.PipelineConfig, error) {
	ret := _m.Called(namespace, name)

	var r0 *v1alpha1.PipelineConfig
	if rf, ok := ret.Get(0).(func(string, string) *v1alpha1.PipelineConfig); ok {
		r0 = rf(namespace, name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1alpha1.PipelineConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPodList provides a mock function with given fields: namespace, opts
func (_m *KubeClient) GetPodList(namespace string, opts v1.ListOptions) (*corev1.PodList, error) {
	ret := _m.Called(namespace, opts)

	var r0 *corev1.PodList
	if rf, ok := ret.Get(0).(func(string, v1.ListOptions) *corev1.PodList); ok {
		r0 = rf(namespace, opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*corev1.PodList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, v1.ListOptions) error); ok {
		r1 = rf(namespace, opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPodNameBylabel provides a mock function with given fields: namespace, label
func (_m *KubeClient) GetPodNameBylabel(namespace string, label string) (string, error) {
	ret := _m.Called(namespace, label)

	var r0 string
	if rf, ok := ret.Get(0).(func(string, string) string); ok {
		r0 = rf(namespace, label)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, label)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// WatchBuildStageStatus provides a mock function with given fields: namespace, label, buildStatus
func (_m *KubeClient) WatchBuildStageStatus(namespace string, label string, buildStatus chan v1alpha1.BuildStatus) error {
	ret := _m.Called(namespace, label, buildStatus)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, chan v1alpha1.BuildStatus) error); ok {
		r0 = rf(namespace, label, buildStatus)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WatchPodStatus provides a mock function with given fields: namespace, label, intervals, podMessage
func (_m *KubeClient) WatchPodStatus(namespace string, label string, intervals int, podMessage chan service.PodMessage) error {
	ret := _m.Called(namespace, label, intervals, podMessage)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int, chan service.PodMessage) error); ok {
		r0 = rf(namespace, label, intervals, podMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
