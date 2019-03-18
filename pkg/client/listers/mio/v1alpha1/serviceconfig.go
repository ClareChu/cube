/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "hidevops.io/mio/pkg/apis/mio/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ServiceConfigLister helps list ServiceConfigs.
type ServiceConfigLister interface {
	// List lists all ServiceConfigs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ServiceConfig, err error)
	// ServiceConfigs returns an object that can list and get ServiceConfigs.
	ServiceConfigs(namespace string) ServiceConfigNamespaceLister
	ServiceConfigListerExpansion
}

// serviceConfigLister implements the ServiceConfigLister interface.
type serviceConfigLister struct {
	indexer cache.Indexer
}

// NewServiceConfigLister returns a new ServiceConfigLister.
func NewServiceConfigLister(indexer cache.Indexer) ServiceConfigLister {
	return &serviceConfigLister{indexer: indexer}
}

// List lists all ServiceConfigs in the indexer.
func (s *serviceConfigLister) List(selector labels.Selector) (ret []*v1alpha1.ServiceConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ServiceConfig))
	})
	return ret, err
}

// ServiceConfigs returns an object that can list and get ServiceConfigs.
func (s *serviceConfigLister) ServiceConfigs(namespace string) ServiceConfigNamespaceLister {
	return serviceConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ServiceConfigNamespaceLister helps list and get ServiceConfigs.
type ServiceConfigNamespaceLister interface {
	// List lists all ServiceConfigs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.ServiceConfig, err error)
	// Get retrieves the ServiceConfig from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.ServiceConfig, error)
	ServiceConfigNamespaceListerExpansion
}

// serviceConfigNamespaceLister implements the ServiceConfigNamespaceLister
// interface.
type serviceConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ServiceConfigs in the indexer for a given namespace.
func (s serviceConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ServiceConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ServiceConfig))
	})
	return ret, err
}

// Get retrieves the ServiceConfig from the indexer for a given namespace and name.
func (s serviceConfigNamespaceLister) Get(name string) (*v1alpha1.ServiceConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("serviceconfig"), name)
	}
	return obj.(*v1alpha1.ServiceConfig), nil
}
