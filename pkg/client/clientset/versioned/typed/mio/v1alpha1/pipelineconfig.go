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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "hidevops.io/mio/pkg/apis/mio/v1alpha1"
	scheme "hidevops.io/mio/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PipelineConfigsGetter has a method to return a PipelineConfigInterface.
// A group's client should implement this interface.
type PipelineConfigsGetter interface {
	PipelineConfigs(namespace string) PipelineConfigInterface
}

// PipelineConfigInterface has methods to work with PipelineConfig resources.
type PipelineConfigInterface interface {
	Create(*v1alpha1.PipelineConfig) (*v1alpha1.PipelineConfig, error)
	Update(*v1alpha1.PipelineConfig) (*v1alpha1.PipelineConfig, error)
	UpdateStatus(*v1alpha1.PipelineConfig) (*v1alpha1.PipelineConfig, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.PipelineConfig, error)
	List(opts v1.ListOptions) (*v1alpha1.PipelineConfigList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PipelineConfig, err error)
	PipelineConfigExpansion
}

// pipelineConfigs implements PipelineConfigInterface
type pipelineConfigs struct {
	client rest.Interface
	ns     string
}

// newPipelineConfigs returns a PipelineConfigs
func newPipelineConfigs(c *MioV1alpha1Client, namespace string) *pipelineConfigs {
	return &pipelineConfigs{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the pipelineConfig, and returns the corresponding pipelineConfig object, and an error if there is any.
func (c *pipelineConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.PipelineConfig, err error) {
	result = &v1alpha1.PipelineConfig{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of PipelineConfigs that match those selectors.
func (c *pipelineConfigs) List(opts v1.ListOptions) (result *v1alpha1.PipelineConfigList, err error) {
	result = &v1alpha1.PipelineConfigList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested pipelineConfigs.
func (c *pipelineConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a pipelineConfig and creates it.  Returns the server's representation of the pipelineConfig, and an error, if there is any.
func (c *pipelineConfigs) Create(pipelineConfig *v1alpha1.PipelineConfig) (result *v1alpha1.PipelineConfig, err error) {
	result = &v1alpha1.PipelineConfig{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		Body(pipelineConfig).
		Do().
		Into(result)
	return
}

// Update takes the representation of a pipelineConfig and updates it. Returns the server's representation of the pipelineConfig, and an error, if there is any.
func (c *pipelineConfigs) Update(pipelineConfig *v1alpha1.PipelineConfig) (result *v1alpha1.PipelineConfig, err error) {
	result = &v1alpha1.PipelineConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		Name(pipelineConfig.Name).
		Body(pipelineConfig).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *pipelineConfigs) UpdateStatus(pipelineConfig *v1alpha1.PipelineConfig) (result *v1alpha1.PipelineConfig, err error) {
	result = &v1alpha1.PipelineConfig{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		Name(pipelineConfig.Name).
		SubResource("status").
		Body(pipelineConfig).
		Do().
		Into(result)
	return
}

// Delete takes name of the pipelineConfig and deletes it. Returns an error if one occurs.
func (c *pipelineConfigs) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *pipelineConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("pipelineconfigs").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched pipelineConfig.
func (c *pipelineConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PipelineConfig, err error) {
	result = &v1alpha1.PipelineConfig{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("pipelineconfigs").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
