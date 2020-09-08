// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/openshift/hive/pkg/apis/hive/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// HiveConfigLister helps list HiveConfigs.
// All objects returned here must be treated as read-only.
type HiveConfigLister interface {
	// List lists all HiveConfigs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.HiveConfig, err error)
	// Get retrieves the HiveConfig from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.HiveConfig, error)
	HiveConfigListerExpansion
}

// hiveConfigLister implements the HiveConfigLister interface.
type hiveConfigLister struct {
	indexer cache.Indexer
}

// NewHiveConfigLister returns a new HiveConfigLister.
func NewHiveConfigLister(indexer cache.Indexer) HiveConfigLister {
	return &hiveConfigLister{indexer: indexer}
}

// List lists all HiveConfigs in the indexer.
func (s *hiveConfigLister) List(selector labels.Selector) (ret []*v1.HiveConfig, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.HiveConfig))
	})
	return ret, err
}

// Get retrieves the HiveConfig from the index for a given name.
func (s *hiveConfigLister) Get(name string) (*v1.HiveConfig, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("hiveconfig"), name)
	}
	return obj.(*v1.HiveConfig), nil
}
