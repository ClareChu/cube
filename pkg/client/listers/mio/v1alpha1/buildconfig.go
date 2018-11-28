/*
Copyright 2018 The Kubernetes Authors.

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

// BuildConfigLister helps list BuildConfigs.
type BuildConfigLister interface {
	// List lists all BuildConfigs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.BuildConfig, err error)
	// BuildConfigs returns an object that can list and get BuildConfigs.
	BuildConfigs(namespace string) BuildConfigNamespaceLister
	BuildConfigListerExpansion
}

// buildConfigLister implements the BuildConfigLister interface.
type buildConfigLister struct {
	indexer cache.Indexer
}

// NewBuildConfigLister returns a new BuildConfigLister.
func NewBuildConfigLister(indexer cache.Indexer) BuildConfigLister {
	return &buildConfigLister{indexer: indexer}
}

// List lists all BuildConfigs in the indexer.
func (s *buildConfigLister) List(selector labels.Selector) (ret []*v1alpha1.BuildConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.BuildConfig))
	})
	return ret, err
}

// BuildConfigs returns an object that can list and get BuildConfigs.
func (s *buildConfigLister) BuildConfigs(namespace string) BuildConfigNamespaceLister {
	return buildConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// BuildConfigNamespaceLister helps list and get BuildConfigs.
type BuildConfigNamespaceLister interface {
	// List lists all BuildConfigs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.BuildConfig, err error)
	// Get retrieves the BuildConfig from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.BuildConfig, error)
	BuildConfigNamespaceListerExpansion
}

// buildConfigNamespaceLister implements the BuildConfigNamespaceLister
// interface.
type buildConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all BuildConfigs in the indexer for a given namespace.
func (s buildConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.BuildConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.BuildConfig))
	})
	return ret, err
}

// Get retrieves the BuildConfig from the indexer for a given namespace and name.
func (s buildConfigNamespaceLister) Get(name string) (*v1alpha1.BuildConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("buildconfig"), name)
	}
	return obj.(*v1alpha1.BuildConfig), nil
}
