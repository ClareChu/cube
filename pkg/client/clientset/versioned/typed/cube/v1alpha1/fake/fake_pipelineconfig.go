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

// FakePipelineConfigs implements PipelineConfigInterface
type FakePipelineConfigs struct {
	Fake *FakeCubeV1alpha1
	ns   string
}

var pipelineconfigsResource = schema.GroupVersionResource{Group: "cube.io", Version: "v1alpha1", Resource: "pipelineconfigs"}

var pipelineconfigsKind = schema.GroupVersionKind{Group: "cube.io", Version: "v1alpha1", Kind: "PipelineConfig"}

// Get takes name of the pipelineConfig, and returns the corresponding pipelineConfig object, and an error if there is any.
func (c *FakePipelineConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.PipelineConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(pipelineconfigsResource, c.ns, name), &v1alpha1.PipelineConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PipelineConfig), err
}

// List takes label and field selectors, and returns the list of PipelineConfigs that match those selectors.
func (c *FakePipelineConfigs) List(opts v1.ListOptions) (result *v1alpha1.PipelineConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(pipelineconfigsResource, pipelineconfigsKind, c.ns, opts), &v1alpha1.PipelineConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PipelineConfigList{}
	for _, item := range obj.(*v1alpha1.PipelineConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested pipelineConfigs.
func (c *FakePipelineConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(pipelineconfigsResource, c.ns, opts))

}

// Create takes the representation of a pipelineConfig and creates it.  Returns the server's representation of the pipelineConfig, and an error, if there is any.
func (c *FakePipelineConfigs) Create(pipelineConfig *v1alpha1.PipelineConfig) (result *v1alpha1.PipelineConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(pipelineconfigsResource, c.ns, pipelineConfig), &v1alpha1.PipelineConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PipelineConfig), err
}

// Update takes the representation of a pipelineConfig and updates it. Returns the server's representation of the pipelineConfig, and an error, if there is any.
func (c *FakePipelineConfigs) Update(pipelineConfig *v1alpha1.PipelineConfig) (result *v1alpha1.PipelineConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(pipelineconfigsResource, c.ns, pipelineConfig), &v1alpha1.PipelineConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PipelineConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePipelineConfigs) UpdateStatus(pipelineConfig *v1alpha1.PipelineConfig) (*v1alpha1.PipelineConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(pipelineconfigsResource, "status", c.ns, pipelineConfig), &v1alpha1.PipelineConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PipelineConfig), err
}

// Delete takes name of the pipelineConfig and deletes it. Returns an error if one occurs.
func (c *FakePipelineConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(pipelineconfigsResource, c.ns, name), &v1alpha1.PipelineConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePipelineConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(pipelineconfigsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.PipelineConfigList{})
	return err
}

// Patch applies the patch and returns the patched pipelineConfig.
func (c *FakePipelineConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.PipelineConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(pipelineconfigsResource, c.ns, name, data, subresources...), &v1alpha1.PipelineConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PipelineConfig), err
}
