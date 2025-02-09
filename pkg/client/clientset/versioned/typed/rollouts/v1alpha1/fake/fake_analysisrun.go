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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/nholuongut/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAnalysisRuns implements AnalysisRunInterface
type FakeAnalysisRuns struct {
	Fake *FakeArgoprojV1alpha1
	ns   string
}

var analysisrunsResource = v1alpha1.SchemeGroupVersion.WithResource("analysisruns")

var analysisrunsKind = v1alpha1.SchemeGroupVersion.WithKind("AnalysisRun")

// Get takes name of the analysisRun, and returns the corresponding analysisRun object, and an error if there is any.
func (c *FakeAnalysisRuns) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AnalysisRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(analysisrunsResource, c.ns, name), &v1alpha1.AnalysisRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AnalysisRun), err
}

// List takes label and field selectors, and returns the list of AnalysisRuns that match those selectors.
func (c *FakeAnalysisRuns) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AnalysisRunList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(analysisrunsResource, analysisrunsKind, c.ns, opts), &v1alpha1.AnalysisRunList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AnalysisRunList{ListMeta: obj.(*v1alpha1.AnalysisRunList).ListMeta}
	for _, item := range obj.(*v1alpha1.AnalysisRunList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested analysisRuns.
func (c *FakeAnalysisRuns) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(analysisrunsResource, c.ns, opts))

}

// Create takes the representation of a analysisRun and creates it.  Returns the server's representation of the analysisRun, and an error, if there is any.
func (c *FakeAnalysisRuns) Create(ctx context.Context, analysisRun *v1alpha1.AnalysisRun, opts v1.CreateOptions) (result *v1alpha1.AnalysisRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(analysisrunsResource, c.ns, analysisRun), &v1alpha1.AnalysisRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AnalysisRun), err
}

// Update takes the representation of a analysisRun and updates it. Returns the server's representation of the analysisRun, and an error, if there is any.
func (c *FakeAnalysisRuns) Update(ctx context.Context, analysisRun *v1alpha1.AnalysisRun, opts v1.UpdateOptions) (result *v1alpha1.AnalysisRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(analysisrunsResource, c.ns, analysisRun), &v1alpha1.AnalysisRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AnalysisRun), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAnalysisRuns) UpdateStatus(ctx context.Context, analysisRun *v1alpha1.AnalysisRun, opts v1.UpdateOptions) (*v1alpha1.AnalysisRun, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(analysisrunsResource, "status", c.ns, analysisRun), &v1alpha1.AnalysisRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AnalysisRun), err
}

// Delete takes name of the analysisRun and deletes it. Returns an error if one occurs.
func (c *FakeAnalysisRuns) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(analysisrunsResource, c.ns, name, opts), &v1alpha1.AnalysisRun{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAnalysisRuns) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(analysisrunsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.AnalysisRunList{})
	return err
}

// Patch applies the patch and returns the patched analysisRun.
func (c *FakeAnalysisRuns) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AnalysisRun, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(analysisrunsResource, c.ns, name, pt, data, subresources...), &v1alpha1.AnalysisRun{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AnalysisRun), err
}
