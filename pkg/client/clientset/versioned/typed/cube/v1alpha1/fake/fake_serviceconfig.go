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

package fake

import (
	v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeServiceConfigs implements ServiceConfigInterface
type FakeServiceConfigs struct {
	Fake *FakeCubeV1alpha1
	ns   string
}

var serviceconfigsResource = schema.GroupVersionResource{Group: "cube.io", Version: "v1alpha1", Resource: "serviceconfigs"}

var serviceconfigsKind = schema.GroupVersionKind{Group: "cube.io", Version: "v1alpha1", Kind: "ServiceConfig"}

// Get takes name of the serviceConfig, and returns the corresponding serviceConfig object, and an error if there is any.
func (c *FakeServiceConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(serviceconfigsResource, c.ns, name), &v1alpha1.ServiceConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceConfig), err
}

// List takes label and field selectors, and returns the list of ServiceConfigs that match those selectors.
func (c *FakeServiceConfigs) List(opts v1.ListOptions) (result *v1alpha1.ServiceConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(serviceconfigsResource, serviceconfigsKind, c.ns, opts), &v1alpha1.ServiceConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ServiceConfigList{}
	for _, item := range obj.(*v1alpha1.ServiceConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested serviceConfigs.
func (c *FakeServiceConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(serviceconfigsResource, c.ns, opts))

}

// Create takes the representation of a serviceConfig and creates it.  Returns the server's representation of the serviceConfig, and an error, if there is any.
func (c *FakeServiceConfigs) Create(serviceConfig *v1alpha1.ServiceConfig) (result *v1alpha1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(serviceconfigsResource, c.ns, serviceConfig), &v1alpha1.ServiceConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceConfig), err
}

// Update takes the representation of a serviceConfig and updates it. Returns the server's representation of the serviceConfig, and an error, if there is any.
func (c *FakeServiceConfigs) Update(serviceConfig *v1alpha1.ServiceConfig) (result *v1alpha1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(serviceconfigsResource, c.ns, serviceConfig), &v1alpha1.ServiceConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeServiceConfigs) UpdateStatus(serviceConfig *v1alpha1.ServiceConfig) (*v1alpha1.ServiceConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(serviceconfigsResource, "status", c.ns, serviceConfig), &v1alpha1.ServiceConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceConfig), err
}

// Delete takes name of the serviceConfig and deletes it. Returns an error if one occurs.
func (c *FakeServiceConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(serviceconfigsResource, c.ns, name), &v1alpha1.ServiceConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeServiceConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(serviceconfigsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.ServiceConfigList{})
	return err
}

// Patch applies the patch and returns the patched serviceConfig.
func (c *FakeServiceConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ServiceConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(serviceconfigsResource, c.ns, name, data, subresources...), &v1alpha1.ServiceConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ServiceConfig), err
}
