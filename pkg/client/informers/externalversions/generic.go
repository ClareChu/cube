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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=cube.io, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("builds"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().Builds().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("buildconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().BuildConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("deployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().Deployments().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("deploymentconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().DeploymentConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("gatewayconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().GatewayConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("imagestreams"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().ImageStreams().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("notifies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().Notifies().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("pipelines"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().Pipelines().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("pipelineconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().PipelineConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("serviceconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().ServiceConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("sourceconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().SourceConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("testconfigs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().TestConfigs().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("testses"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Cube().V1alpha1().Testses().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
