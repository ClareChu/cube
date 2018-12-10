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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	miov1alpha1 "hidevops.io/mio/pkg/apis/mio/v1alpha1"
	versioned "hidevops.io/mio/pkg/client/clientset/versioned"
	internalinterfaces "hidevops.io/mio/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "hidevops.io/mio/pkg/client/listers/mio/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// TestsInformer provides access to a shared informer and lister for
// Testses.
type TestsInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.TestsLister
}

type testsInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewTestsInformer constructs a new informer for Tests type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewTestsInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredTestsInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredTestsInformer constructs a new informer for Tests type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredTestsInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MioV1alpha1().Testses(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.MioV1alpha1().Testses(namespace).Watch(options)
			},
		},
		&miov1alpha1.Tests{},
		resyncPeriod,
		indexers,
	)
}

func (f *testsInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredTestsInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *testsInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&miov1alpha1.Tests{}, f.defaultInformer)
}

func (f *testsInformer) Lister() v1alpha1.TestsLister {
	return v1alpha1.NewTestsLister(f.Informer().GetIndexer())
}
