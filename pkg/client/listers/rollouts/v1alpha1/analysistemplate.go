/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/nholuongut/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// AnalysisTemplateLister helps list AnalysisTemplates.
// All objects returned here must be treated as read-only.
type AnalysisTemplateLister interface {
	// List lists all AnalysisTemplates in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.AnalysisTemplate, err error)
	// AnalysisTemplates returns an object that can list and get AnalysisTemplates.
	AnalysisTemplates(namespace string) AnalysisTemplateNamespaceLister
	AnalysisTemplateListerExpansion
}

// analysisTemplateLister implements the AnalysisTemplateLister interface.
type analysisTemplateLister struct {
	indexer cache.Indexer
}

// NewAnalysisTemplateLister returns a new AnalysisTemplateLister.
func NewAnalysisTemplateLister(indexer cache.Indexer) AnalysisTemplateLister {
	return &analysisTemplateLister{indexer: indexer}
}

// List lists all AnalysisTemplates in the indexer.
func (s *analysisTemplateLister) List(selector labels.Selector) (ret []*v1alpha1.AnalysisTemplate, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AnalysisTemplate))
	})
	return ret, err
}

// AnalysisTemplates returns an object that can list and get AnalysisTemplates.
func (s *analysisTemplateLister) AnalysisTemplates(namespace string) AnalysisTemplateNamespaceLister {
	return analysisTemplateNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// AnalysisTemplateNamespaceLister helps list and get AnalysisTemplates.
// All objects returned here must be treated as read-only.
type AnalysisTemplateNamespaceLister interface {
	// List lists all AnalysisTemplates in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.AnalysisTemplate, err error)
	// Get retrieves the AnalysisTemplate from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.AnalysisTemplate, error)
	AnalysisTemplateNamespaceListerExpansion
}

// analysisTemplateNamespaceLister implements the AnalysisTemplateNamespaceLister
// interface.
type analysisTemplateNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all AnalysisTemplates in the indexer for a given namespace.
func (s analysisTemplateNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.AnalysisTemplate, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.AnalysisTemplate))
	})
	return ret, err
}

// Get retrieves the AnalysisTemplate from the indexer for a given namespace and name.
func (s analysisTemplateNamespaceLister) Get(name string) (*v1alpha1.AnalysisTemplate, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("analysistemplate"), name)
	}
	return obj.(*v1alpha1.AnalysisTemplate), nil
}
