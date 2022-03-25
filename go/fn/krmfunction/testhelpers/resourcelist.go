// // Copyright 2022 Google LLC
// //
// // Licensed under the Apache License, Version 2.0 (the "License");
// // you may not use this file except in compliance with the License.
// // You may obtain a copy of the License at
// //
// //      http://www.apache.org/licenses/LICENSE-2.0
// //
// // Unless required by applicable law or agreed to in writing, software
// // distributed under the License is distributed on an "AS IS" BASIS,
// // WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// // See the License for the specific language governing permissions and
// // limitations under the License.

package testhelpers

// import (
// 	"fmt"
// 	"strings"

// 	"github.com/GoogleContainerTools/kpt/porch/functions/pkg/yo"
// )

// type ResourceList struct {
// 	doc         *yo.Doc
// 	kubeObjects []*yo.KubeObject
// }

// func NewResourceList(s string) (*ResourceList, error) {
// 	doc, err := yo.Parse([]byte(s))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse: %w", err)
// 	}
// 	kubeObjects := yo.ExtractKubeObjects(doc)

// 	return &ResourceList{
// 		doc:         doc,
// 		kubeObjects: kubeObjects,
// 	}, nil
// }

// func (r *ResourceList) Debug() string {
// 	var info []string
// 	for _, obj := range r.kubeObjects {
// 		info = append(info, fmt.Sprintf("%s::%s/%s", obj.GetKind(), obj.GetNamespace(), obj.GetName()))
// 	}
// 	return strings.Join(info, ";")
// 	// b, err := r.doc.ToYAML()
// 	// if err != nil {
// 	// 	return fmt.Sprintf("<ToYAML failed: %v>", err)
// 	// }
// 	// return string(b)
// }

// func (r *ResourceList) KubernetesObjects() *KubeObjectList {
// 	return &KubeObjectList{
// 		objects: r.kubeObjects,
// 	}
// }

// func (r *ResourceList) AddKubeObject(apiVersion, kind string) (*yo.KubeObject, error) {
// 	objects := r.doc.Objects()

// 	if len(objects) != 1 {
// 		return nil, fmt.Errorf("expected exactly one ResourceList object")
// 	}
// 	rl := yo.AsKubeObject(objects[0])
// 	if rl.GetKind() != "ResourceList" {
// 		return nil, fmt.Errorf("expected object to be of kind ResourceList")
// 	}
// 	if rl.GetAPIVersion() != "config.kubernetes.io/v1alpha1" {
// 		return nil, fmt.Errorf("expected object to be of apiVersion config.kubernetes.io/v1alpha1")
// 	}

// 	obj := yo.NewMap()
// 	ko := yo.AsKubeObject(obj)
// 	ko.SetAPIVersion(apiVersion)
// 	ko.SetKind(kind)

// 	items := rl.UpsertSlice("items")

// 	items.Add(ko.Node())

// 	return ko, nil
// }

// func (r *ResourceList) Encode() ([]byte, error) {
// 	return r.doc.ToYAML()
// }

// func (r *ResourceList) Decode(b []byte) error {
// 	rl, err := NewResourceList(string(b))
// 	if err != nil {
// 		return err
// 	}
// 	*r = *rl
// 	return nil
// }

// type KubeObjectList struct {
// 	objects []*yo.KubeObject
// }

// func (l *KubeObjectList) OfKind(apiVersion, kind string) *KubeObjectList {
// 	ret := &KubeObjectList{}
// 	for _, obj := range l.objects {
// 		if obj.GetAPIVersion() == apiVersion && obj.GetKind() == kind {
// 			ret.objects = append(ret.objects, obj)
// 		}
// 	}
// 	return ret
// }

// func (l *KubeObjectList) ExactlyOne() (*yo.KubeObject, error) {
// 	if len(l.objects) == 0 {
// 		return nil, fmt.Errorf("found no matching objects, expected exactly one")
// 	}
// 	if len(l.objects) > 1 {
// 		return nil, fmt.Errorf("found multiple (%d) matching objects, expected exactly one", len(l.objects))
// 	}
// 	return l.objects[0], nil
// }

// func (l *KubeObjectList) ForEach(fn func(*yo.KubeObject) error) error {
// 	for _, obj := range l.objects {
// 		if err := fn(obj); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
