// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package krmfunctiontests

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildResourceList(t *testing.T) {
	items := `
kind: ConfigMap
apiVersion: v1
metadata:
  name: foo
`

	config := `
kind: ConfigMap
apiVersion: v1
data:
  paramkey: paramvalue
`

	itemsObjects, err := parseObjects([]byte(items))
	if err != nil {
		t.Fatalf("parseObjects failed: %v", err)
	}
	configObjects, err := parseObjects([]byte(config))
	if err != nil {
		t.Fatalf("parseObjects failed: %v", err)
	}
	got, err := buildResourceList(itemsObjects, configObjects)
	if err != nil {
		t.Fatalf("buildResourceList failed: %v", err)
	}

	want := `
apiVersion: config.kubernetes.io/v1alpha1
kind: ResourceList
items:
- kind: ConfigMap
  apiVersion: v1
  metadata:
    name: foo
functionConfig:
  kind: ConfigMap
  apiVersion: v1
  data:
    paramkey: paramvalue
`

	if diff := cmp.Diff(strings.TrimSpace(string(got)), strings.TrimSpace(want)); diff != "" {
		t.Logf("got: %s", got)

		t.Fatalf("unexpected diff: %s", diff)
	}
}
