/*
Copyright 2018 caicloud authors. All rights reserved.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	kubernetes "github.com/caicloud/clientset/kubernetes"
	v1alpha1 "github.com/caicloud/clientset/listers/config/v1alpha1"
	configv1alpha1 "github.com/caicloud/clientset/pkg/apis/config/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	internalinterfaces "k8s.io/client-go/informers/internalinterfaces"
	clientgokubernetes "k8s.io/client-go/kubernetes"
	cache "k8s.io/client-go/tools/cache"
)

// ConfigReferenceInformer provides access to a shared informer and lister for
// ConfigReferences.
type ConfigReferenceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ConfigReferenceLister
}

type configReferenceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewConfigReferenceInformer constructs a new informer for ConfigReference type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewConfigReferenceInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredConfigReferenceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredConfigReferenceInformer constructs a new informer for ConfigReference type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredConfigReferenceInformer(client kubernetes.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1alpha1().ConfigReferences(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ConfigV1alpha1().ConfigReferences(namespace).Watch(options)
			},
		},
		&configv1alpha1.ConfigReference{},
		resyncPeriod,
		indexers,
	)
}

func (f *configReferenceInformer) defaultInformer(client clientgokubernetes.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredConfigReferenceInformer(client.(kubernetes.Interface), f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *configReferenceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&configv1alpha1.ConfigReference{}, f.defaultInformer)
}

func (f *configReferenceInformer) Lister() v1alpha1.ConfigReferenceLister {
	return v1alpha1.NewConfigReferenceLister(f.Informer().GetIndexer())
}
