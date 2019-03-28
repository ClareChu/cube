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

// FakeGatewayConfigs implements GatewayConfigInterface
type FakeGatewayConfigs struct {
	Fake *FakeCubeV1alpha1
	ns   string
}

var gatewayconfigsResource = schema.GroupVersionResource{Group: "cube.io", Version: "v1alpha1", Resource: "gatewayconfigs"}

var gatewayconfigsKind = schema.GroupVersionKind{Group: "cube.io", Version: "v1alpha1", Kind: "GatewayConfig"}

// Get takes name of the gatewayConfig, and returns the corresponding gatewayConfig object, and an error if there is any.
func (c *FakeGatewayConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.GatewayConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(gatewayconfigsResource, c.ns, name), &v1alpha1.GatewayConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GatewayConfig), err
}

// List takes label and field selectors, and returns the list of GatewayConfigs that match those selectors.
func (c *FakeGatewayConfigs) List(opts v1.ListOptions) (result *v1alpha1.GatewayConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(gatewayconfigsResource, gatewayconfigsKind, c.ns, opts), &v1alpha1.GatewayConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.GatewayConfigList{}
	for _, item := range obj.(*v1alpha1.GatewayConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested gatewayConfigs.
func (c *FakeGatewayConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(gatewayconfigsResource, c.ns, opts))

}

// Create takes the representation of a gatewayConfig and creates it.  Returns the server's representation of the gatewayConfig, and an error, if there is any.
func (c *FakeGatewayConfigs) Create(gatewayConfig *v1alpha1.GatewayConfig) (result *v1alpha1.GatewayConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(gatewayconfigsResource, c.ns, gatewayConfig), &v1alpha1.GatewayConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GatewayConfig), err
}

// Update takes the representation of a gatewayConfig and updates it. Returns the server's representation of the gatewayConfig, and an error, if there is any.
func (c *FakeGatewayConfigs) Update(gatewayConfig *v1alpha1.GatewayConfig) (result *v1alpha1.GatewayConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(gatewayconfigsResource, c.ns, gatewayConfig), &v1alpha1.GatewayConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GatewayConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeGatewayConfigs) UpdateStatus(gatewayConfig *v1alpha1.GatewayConfig) (*v1alpha1.GatewayConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(gatewayconfigsResource, "status", c.ns, gatewayConfig), &v1alpha1.GatewayConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GatewayConfig), err
}

// Delete takes name of the gatewayConfig and deletes it. Returns an error if one occurs.
func (c *FakeGatewayConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(gatewayconfigsResource, c.ns, name), &v1alpha1.GatewayConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeGatewayConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(gatewayconfigsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.GatewayConfigList{})
	return err
}

// Patch applies the patch and returns the patched gatewayConfig.
func (c *FakeGatewayConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GatewayConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(gatewayconfigsResource, c.ns, name, data, subresources...), &v1alpha1.GatewayConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.GatewayConfig), err
}