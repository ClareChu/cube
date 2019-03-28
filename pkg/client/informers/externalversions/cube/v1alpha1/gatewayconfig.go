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

package v1alpha1

import (
	time "time"

	cubev1alpha1 "hidevops.io/cube/pkg/apis/cube/v1alpha1"
	versioned "hidevops.io/cube/pkg/client/clientset/versioned"
	internalinterfaces "hidevops.io/cube/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "hidevops.io/cube/pkg/client/listers/cube/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// GatewayConfigInformer provides access to a shared informer and lister for
// GatewayConfigs.
type GatewayConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.GatewayConfigLister
}

type gatewayConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewGatewayConfigInformer constructs a new informer for GatewayConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewGatewayConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredGatewayConfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredGatewayConfigInformer constructs a new informer for GatewayConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredGatewayConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CubeV1alpha1().GatewayConfigs(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.CubeV1alpha1().GatewayConfigs(namespace).Watch(options)
			},
		},
		&cubev1alpha1.GatewayConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *gatewayConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredGatewayConfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *gatewayConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&cubev1alpha1.GatewayConfig{}, f.defaultInformer)
}

func (f *gatewayConfigInformer) Lister() v1alpha1.GatewayConfigLister {
	return v1alpha1.NewGatewayConfigLister(f.Informer().GetIndexer())
}