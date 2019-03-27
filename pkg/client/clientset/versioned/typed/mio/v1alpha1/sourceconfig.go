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

// SourceConfigsGetter has a method to return a SourceConfigInterface.
// A group's client should implement this interface.
type SourceConfigsGetter interface {
	SourceConfigs(namespace string) SourceConfigInterface
}

// SourceConfigInterface has methods to work with SourceConfig resources.
type SourceConfigInterface interface {
	Create(*v1alpha1.SourceConfig) (*v1alpha1.SourceConfig, error)
	Update(*v1alpha1.SourceConfig) (*v1alpha1.SourceConfig, error)
	UpdateStatus(*v1alpha1.SourceConfig) (*v1alpha1.SourceConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.SourceConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.SourceConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SourceConfig, err error)
	SourceConfigExpansion
}

// sourceConfigs implements SourceConfigInterface
type sourceConfigs struct {
	client rest.Interface
	ns     string
}

// newSourceConfigs returns a SourceConfigs
func newSourceConfigs(c *CubeV1alpha1Client, namespace string) *sourceConfigs {
	return &sourceConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the sourceConfig, and returns the corresponding sourceConfig object, and an error if there is any.
func (c *sourceConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.SourceConfig, err error) {
	result = &v1alpha1.SourceConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sourceconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SourceConfigs that match those selectors.
func (c *sourceConfigs) List(opts v1.ListOptions) (result *v1alpha1.SourceConfigList, err error) {
	result = &v1alpha1.SourceConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("sourceconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested sourceConfigs.
func (c *sourceConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("sourceconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a sourceConfig and creates it.  Returns the server's representation of the sourceConfig, and an error, if there is any.
func (c *sourceConfigs) Create(sourceConfig *v1alpha1.SourceConfig) (result *v1alpha1.SourceConfig, err error) {
	result = &v1alpha1.SourceConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("sourceconfigs").
		Body(sourceConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a sourceConfig and updates it. Returns the server's representation of the sourceConfig, and an error, if there is any.
func (c *sourceConfigs) Update(sourceConfig *v1alpha1.SourceConfig) (result *v1alpha1.SourceConfig, err error) {
	result = &v1alpha1.SourceConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sourceconfigs").
		Name(sourceConfig.Name).
		Body(sourceConfig).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *sourceConfigs) UpdateStatus(sourceConfig *v1alpha1.SourceConfig) (result *v1alpha1.SourceConfig, err error) {
	result = &v1alpha1.SourceConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("sourceconfigs").
		Name(sourceConfig.Name).
		SubResource("status").
		Body(sourceConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the sourceConfig and deletes it. Returns an error if one occurs.
func (c *sourceConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sourceconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *sourceConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("sourceconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched sourceConfig.
func (c *sourceConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SourceConfig, err error) {
	result = &v1alpha1.SourceConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("sourceconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
