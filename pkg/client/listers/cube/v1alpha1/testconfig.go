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
	v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// TestConfigLister helps list TestConfigs.
type TestConfigLister interface {
	// List lists all TestConfigs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.TestConfig, err error)
	// TestConfigs returns an object that can list and get TestConfigs.
	TestConfigs(namespace string) TestConfigNamespaceLister
	TestConfigListerExpansion
}

// testConfigLister implements the TestConfigLister interface.
type testConfigLister struct {
	indexer cache.Indexer
}

// NewTestConfigLister returns a new TestConfigLister.
func NewTestConfigLister(indexer cache.Indexer) TestConfigLister {
	return &testConfigLister{indexer: indexer}
}

// List lists all TestConfigs in the indexer.
func (s *testConfigLister) List(selector labels.Selector) (ret []*v1alpha1.TestConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TestConfig))
	})
	return ret, err
}

// TestConfigs returns an object that can list and get TestConfigs.
func (s *testConfigLister) TestConfigs(namespace string) TestConfigNamespaceLister {
	return testConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// TestConfigNamespaceLister helps list and get TestConfigs.
type TestConfigNamespaceLister interface {
	// List lists all TestConfigs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.TestConfig, err error)
	// Get retrieves the TestConfig from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.TestConfig, error)
	TestConfigNamespaceListerExpansion
}

// testConfigNamespaceLister implements the TestConfigNamespaceLister
// interface.
type testConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all TestConfigs in the indexer for a given namespace.
func (s testConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.TestConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.TestConfig))
	})
	return ret, err
}

// Get retrieves the TestConfig from the indexer for a given namespace and name.
func (s testConfigNamespaceLister) Get(name string) (*v1alpha1.TestConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("testconfig"), name)
	}
	return obj.(*v1alpha1.TestConfig), nil
}
