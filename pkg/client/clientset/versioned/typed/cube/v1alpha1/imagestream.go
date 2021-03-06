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
	v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	scheme "hidevops.io/cube/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ImageStreamsGetter has a method to return a ImageStreamInterface.
// A group's client should implement this interface.
type ImageStreamsGetter interface {
	ImageStreams(namespace string) ImageStreamInterface
}

// ImageStreamInterface has methods to work with ImageStream resources.
type ImageStreamInterface interface {
	Create(*v1alpha1.ImageStream) (*v1alpha1.ImageStream, error)
	Update(*v1alpha1.ImageStream) (*v1alpha1.ImageStream, error)
	UpdateStatus(*v1alpha1.ImageStream) (*v1alpha1.ImageStream, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.ImageStream, error)
	List(opts v1.ListOptions) (*v1alpha1.ImageStreamList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ImageStream, err error)
	ImageStreamExpansion
}

// imageStreams implements ImageStreamInterface
type imageStreams struct {
	client rest.Interface
	ns     string
}

// newImageStreams returns a ImageStreams
func newImageStreams(c *CubeV1alpha1Client, namespace string) *imageStreams {
	return &imageStreams{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the imageStream, and returns the corresponding imageStream object, and an error if there is any.
func (c *imageStreams) Get(name string, options v1.GetOptions) (result *v1alpha1.ImageStream, err error) {
	result = &v1alpha1.ImageStream{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("imagestreams").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ImageStreams that match those selectors.
func (c *imageStreams) List(opts v1.ListOptions) (result *v1alpha1.ImageStreamList, err error) {
	result = &v1alpha1.ImageStreamList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("imagestreams").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested imageStreams.
func (c *imageStreams) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("imagestreams").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a imageStream and creates it.  Returns the server's representation of the imageStream, and an error, if there is any.
func (c *imageStreams) Create(imageStream *v1alpha1.ImageStream) (result *v1alpha1.ImageStream, err error) {
	result = &v1alpha1.ImageStream{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("imagestreams").
		Body(imageStream).
		Do().
		Into(result)
	return
}

// Update takes the representation of a imageStream and updates it. Returns the server's representation of the imageStream, and an error, if there is any.
func (c *imageStreams) Update(imageStream *v1alpha1.ImageStream) (result *v1alpha1.ImageStream, err error) {
	result = &v1alpha1.ImageStream{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("imagestreams").
		Name(imageStream.Name).
		Body(imageStream).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *imageStreams) UpdateStatus(imageStream *v1alpha1.ImageStream) (result *v1alpha1.ImageStream, err error) {
	result = &v1alpha1.ImageStream{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("imagestreams").
		Name(imageStream.Name).
		SubResource("status").
		Body(imageStream).
		Do().
		Into(result)
	return
}

// Delete takes name of the imageStream and deletes it. Returns an error if one occurs.
func (c *imageStreams) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("imagestreams").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *imageStreams) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("imagestreams").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched imageStream.
func (c *imageStreams) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ImageStream, err error) {
	result = &v1alpha1.ImageStream{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("imagestreams").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
