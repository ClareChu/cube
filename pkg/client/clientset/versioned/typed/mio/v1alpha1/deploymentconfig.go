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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"hidevops.io/cube/pkg/apis/cube/v1alpha1"
	"hidevops.io/cube/pkg/client/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

// DeploymentConfigsGetter has a method to return a DeploymentConfigInterface.
// A group's client should implement this interface.
type DeploymentConfigsGetter interface {
	DeploymentConfigs(namespace string) DeploymentConfigInterface
}

// DeploymentConfigInterface has methods to work with DeploymentConfig resources.
type DeploymentConfigInterface interface {
	Create(*v1alpha1.DeploymentConfig) (*v1alpha1.DeploymentConfig, error)
	Update(*v1alpha1.DeploymentConfig) (*v1alpha1.DeploymentConfig, error)
	UpdateStatus(*v1alpha1.DeploymentConfig) (*v1alpha1.DeploymentConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.DeploymentConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.DeploymentConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DeploymentConfig, err error)
	DeploymentConfigExpansion
}

// deploymentConfigs implements DeploymentConfigInterface
type deploymentConfigs struct {
	client rest.Interface
	ns     string
}

// newDeploymentConfigs returns a DeploymentConfigs
func newDeploymentConfigs(c *CubeV1alpha1Client, namespace string) *deploymentConfigs {
	return &deploymentConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the deploymentConfig, and returns the corresponding deploymentConfig object, and an error if there is any.
func (c *deploymentConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.DeploymentConfig, err error) {
	result = &v1alpha1.DeploymentConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DeploymentConfigs that match those selectors.
func (c *deploymentConfigs) List(opts v1.ListOptions) (result *v1alpha1.DeploymentConfigList, err error) {
	result = &v1alpha1.DeploymentConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested deploymentConfigs.
func (c *deploymentConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a deploymentConfig and creates it.  Returns the server's representation of the deploymentConfig, and an error, if there is any.
func (c *deploymentConfigs) Create(deploymentConfig *v1alpha1.DeploymentConfig) (result *v1alpha1.DeploymentConfig, err error) {
	result = &v1alpha1.DeploymentConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		Body(deploymentConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a deploymentConfig and updates it. Returns the server's representation of the deploymentConfig, and an error, if there is any.
func (c *deploymentConfigs) Update(deploymentConfig *v1alpha1.DeploymentConfig) (result *v1alpha1.DeploymentConfig, err error) {
	result = &v1alpha1.DeploymentConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		Name(deploymentConfig.Name).
		Body(deploymentConfig).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *deploymentConfigs) UpdateStatus(deploymentConfig *v1alpha1.DeploymentConfig) (result *v1alpha1.DeploymentConfig, err error) {
	result = &v1alpha1.DeploymentConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		Name(deploymentConfig.Name).
		SubResource("status").
		Body(deploymentConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the deploymentConfig and deletes it. Returns an error if one occurs.
func (c *deploymentConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *deploymentConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("deploymentconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched deploymentConfig.
func (c *deploymentConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.DeploymentConfig, err error) {
	result = &v1alpha1.DeploymentConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("deploymentconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
