/*
Copyright 2018 caicloud authors. All rights reserved.
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/caicloud/clientset/pkg/apis/tenant/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PartitionLister helps list Partitions.
type PartitionLister interface {
	// List lists all Partitions in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Partition, err error)
	// Get retrieves the Partition from the index for a given name.
	Get(name string) (*v1alpha1.Partition, error)
	PartitionListerExpansion
}

// partitionLister implements the PartitionLister interface.
type partitionLister struct {
	indexer cache.Indexer
}

// NewPartitionLister returns a new PartitionLister.
func NewPartitionLister(indexer cache.Indexer) PartitionLister {
	return &partitionLister{indexer: indexer}
}

// List lists all Partitions in the indexer.
func (s *partitionLister) List(selector labels.Selector) (ret []*v1alpha1.Partition, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Partition))
	})
	return ret, err
}

// Get retrieves the Partition from the index for a given name.
func (s *partitionLister) Get(name string) (*v1alpha1.Partition, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("partition"), name)
	}
	return obj.(*v1alpha1.Partition), nil
}
