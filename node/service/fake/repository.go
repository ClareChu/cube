package fake

import (
	"github.com/stretchr/testify/mock"
	cubev1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type Repository struct {
	mock.Mock
}

func (e *Repository) Create(build *cubev1alpha1.Build) (config *cubev1alpha1.Build, err error) {
	args := e.Called(build)
	return args[0].(*cubev1alpha1.Build), args.Error(1)
}

func (e *Repository) Get(name, namespace string) (config *cubev1alpha1.Build, err error) {
	args := e.Called(name, namespace)
	return args[0].(*cubev1alpha1.Build), args.Error(1)
}

func (e *Repository) Watch(listOptions v1.ListOptions, namespace, name string) (watch.Interface, error) {
	args := e.Called(nil, name, namespace)
	return args[0].(watch.Interface), args.Error(1)
}

func (e *Repository) Delete(name, namespace string) error {
	args := e.Called(name, namespace)
	return args.Error(1)
}

func (e *Repository) Update(name, namespace string, config *cubev1alpha1.Build) (*cubev1alpha1.Build, error) {
	args := e.Called(name, namespace, config)
	return args[0].(*cubev1alpha1.Build), args.Error(1)
}

func (e *Repository) UpdateBuildStatus(namespace, name, eventType, status string) (*cubev1alpha1.Build, error) {
	args := e.Called(name, namespace, eventType, status)
	return args[0].(*cubev1alpha1.Build), args.Error(1)
}
