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
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn/krmfunction"
	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn/krmfunction/testhelpers"
)

func RunFunctionObjectTests(t *testing.T, testdata string, fnPrototype krmfunction.KRMFunctionObject) {
	dirs, err := ioutil.ReadDir(testdata)
	if err != nil {
		t.Fatalf("failed to read directory %q: %v", testdata, err)
	}

	for _, d := range dirs {
		dir := filepath.Join(testdata, d.Name())
		if !d.IsDir() {
			t.Errorf("expected directory, found %s", dir)
			continue
		}

		t.Run(dir, func(t *testing.T) {
			items := mustParseFile(t, filepath.Join(dir, "in.yaml"))
			config := mustParseFile(t, filepath.Join(dir, "config.yaml"))

			var functionConfig *fn.KubeObject
			if len(config) == 0 {
				functionConfig = nil
			} else if len(config) == 1 {
				functionConfig = config[0]
			} else {
				t.Fatalf("found multiple config objects in %s", filepath.Join(dir, "config.yaml"))
			}

			rl, err := buildResourceList(items, functionConfig)
			if err != nil {
				t.Fatalf("failed to build resource list: %v", err)
			}

			// t.Logf("stdin: %v", string(rl))

			fnType := reflect.TypeOf(fnPrototype).Elem()

			fnValue := reflect.New(fnType)
			fnObject := fnValue.Interface().(krmfunction.KRMFunctionObject)

			processor := krmfunction.ConvertObjectToResourceListProcessor(fnObject)

			if err := processor.Process(rl); err != nil {
				t.Fatalf("Run failed unexpectedly: %v", err)
			}

			rlYAML, err := rl.ToYAML()
			if err != nil {
				t.Fatalf("failed to convert resource list to yaml: %v", err)
			}

			p := filepath.Join(dir, "expected.yaml")
			testhelpers.CompareGoldenFile(t, p, rlYAML)
		})
	}
}

// func RunTests(t *testing.T, testdata string, function krmfunction.KRMFunction) {
// 	dirs, err := ioutil.ReadDir(testdata)
// 	if err != nil {
// 		t.Fatalf("failed to read directory %q: %v", testdata, err)
// 	}

// 	for _, d := range dirs {
// 		dir := filepath.Join(testdata, d.Name())
// 		if !d.IsDir() {
// 			t.Errorf("expected directory, found %s", dir)
// 			continue
// 		}

// 		t.Run(dir, func(t *testing.T) {
// 			items := mustParseFile(t, filepath.Join(dir, "in.yaml"))
// 			config := mustParseFile(t, filepath.Join(dir, "config.yaml"))

// 			var functionConfig *fn.KubeObject
// 			if len(config) == 0 {
// 				functionConfig = nil
// 			} else if len(config) == 1 {
// 				functionConfig = config[0]
// 			} else {
// 				t.Fatalf("found multiple config objects in %s", filepath.Join(dir, "config.yaml"))
// 			}

// 			rl, err := buildResourceList(items, functionConfig)
// 			if err != nil {
// 				t.Fatalf("failed to build resource list: %v", err)
// 			}

// 			t.Logf("stdin: %v", string(rl))

// 			krmContext := &fn.Context{Context: context.Background()}
// 			stdout, err := krmfunction.Run(krmContext, function, rl)
// 			if err != nil {
// 				t.Fatalf("Run failed unexpectedly: %v", err)
// 			}

// 			p := filepath.Join(dir, "expected.yaml")
// 			testhelpers.CompareGoldenFile(t, p, stdout)
// 		})
// 	}

// }

func mustParseFile(t *testing.T, p string) []*fn.KubeObject {
	b := testhelpers.MustReadFile(t, p)
	objects, err := fn.ParseKubeObjects(b)
	if err != nil {
		t.Fatalf("failed to parse objects from file %q: %v", p, err)
	}
	return objects
}

func buildResourceList(items []*fn.KubeObject, functionConfig *fn.KubeObject) (*fn.ResourceList, error) {
	rl := &fn.ResourceList{}
	rl.Items = items
	rl.FunctionConfig = functionConfig

	return rl, nil
	// doc := fn.NewDoc()
	// rl := doc.NewMap()
	// rlObj := fn.AsKubeObject(rl)
	// rlObj.SetAPIVersion("config.kubernetes.io/v1alpha1")
	// rlObj.SetKind("ResourceList")

	// itemSlice := rlObj.UpsertSlice("items")
	// for _, item := range items {
	// 	itemSlice.Add(item.Node())
	// }

	// if len(config) > 1 {
	// 	return nil, fmt.Errorf("expected exactly one config object, got %d", len(config))
	// }

	// if len(config) == 1 {
	// 	rlObj.Node().Set("functionConfig", config[0].Node().Node())
	// }

	// y, err := doc.ToYAML()
	// y, err := rl.ToYAML()
	// if err != nil {
	// 	return nil, err
	// }
	// return y, nil
}
