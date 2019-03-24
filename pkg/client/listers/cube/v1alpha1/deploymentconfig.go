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

// DeploymentConfigLister helps list DeploymentConfigs.
type DeploymentConfigLister interface {
	// List lists all DeploymentConfigs in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.DeploymentConfig, err error)
	// DeploymentConfigs returns an object that can list and get DeploymentConfigs.
	DeploymentConfigs(namespace string) DeploymentConfigNamespaceLister
	DeploymentConfigListerExpansion
}

// deploymentConfigLister implements the DeploymentConfigLister interface.
type deploymentConfigLister struct {
	indexer cache.Indexer
}

// NewDeploymentConfigLister returns a new DeploymentConfigLister.
func NewDeploymentConfigLister(indexer cache.Indexer) DeploymentConfigLister {
	return &deploymentConfigLister{indexer: indexer}
}

// List lists all DeploymentConfigs in the indexer.
func (s *deploymentConfigLister) List(selector labels.Selector) (ret []*v1alpha1.DeploymentConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DeploymentConfig))
	})
	return ret, err
}

// DeploymentConfigs returns an object that can list and get DeploymentConfigs.
func (s *deploymentConfigLister) DeploymentConfigs(namespace string) DeploymentConfigNamespaceLister {
	return deploymentConfigNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// DeploymentConfigNamespaceLister helps list and get DeploymentConfigs.
type DeploymentConfigNamespaceLister interface {
	// List lists all DeploymentConfigs in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.DeploymentConfig, err error)
	// Get retrieves the DeploymentConfig from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.DeploymentConfig, error)
	DeploymentConfigNamespaceListerExpansion
}

// deploymentConfigNamespaceLister implements the DeploymentConfigNamespaceLister
// interface.
type deploymentConfigNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all DeploymentConfigs in the indexer for a given namespace.
func (s deploymentConfigNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.DeploymentConfig, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.DeploymentConfig))
	})
	return ret, err
}

// Get retrieves the DeploymentConfig from the indexer for a given namespace and name.
func (s deploymentConfigNamespaceLister) Get(name string) (*v1alpha1.DeploymentConfig, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("deploymentconfig"), name)
	}
	return obj.(*v1alpha1.DeploymentConfig), nil
}
