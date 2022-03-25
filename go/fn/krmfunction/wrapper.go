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

package krmfunction

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"os"

// 	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
// 	"k8s.io/klog/v2"
// )

// type KRMFunction func(ctx context.Context, in []*fn.KubeObject, config map[string]string) ([]*fn.KubeObject, error)

// // RunAsMain should be called from main
// func RunAsMain(fn KRMFunction) {
// 	klog.InitFlags(nil)
// 	flag.Parse()

// 	in, err := ioutil.ReadAll(os.Stdin)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "failed to read from stdin: %v\n", err)
// 		os.Exit(1)
// 	}

// 	ctx := context.Background()

// 	krmContext := &Context{Context: ctx}
// 	stdout, err := Run(krmContext, fn, in)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "%v\n", err)
// 		os.Exit(1)
// 	}

// 	if _, err := os.Stdout.Write(stdout); err != nil {
// 		fmt.Fprintf(os.Stderr, "error writing response to stdout: %v\n", err)
// 		os.Exit(1)
// 	}

// 	if len(krmContext.Results) != 0 {
// 		// Based on https://github.com/GoogleContainerTools/kpt/issues/2536, my interpretation is that we return 0 if we ran without internal problems
// 		// But we also output structured error information on errors
// 		os.Exit(0)
// 	} else {
// 		os.Exit(0)
// 	}
// }

// // Run evaluates the function
// func Run(krmContext *Context, function KRMFunction, stdin []byte) ([]byte, error) {
// 	resourceList, err := fn.ParseResourceList(stdin)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse stdin: %w", err)
// 	}

// 	objects := doc.Objects()
// 	if len(objects) != 1 {
// 		return nil, fmt.Errorf("unexpected input, expected exactly one object, got %d", len(objects))
// 	}

// 	obj := fn.AsKubeObject(objects[0])
// 	if obj.GetKind() != "ResourceList" {
// 		return nil, fmt.Errorf("input was of unexpected kind %q; expected ResourceList", obj.GetKind())
// 	}

// 	var items []*fn.KubeObject
// 	itemsSlice := obj.UpsertSlice("items")

// 	itemsObjects := itemsSlice.Objects()

// 	for _, itemsObject := range itemsObjects {
// 		items = append(items, fn.AsKubeObject(itemsObject))
// 	}

// 	functionConfigMap, found, err := obj.GetMap("functionConfig")
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting functionConfig: %w", err)
// 	}
// 	functionConfig := map[string]string{}
// 	if found {
// 		functionConfigObject := fn.AsKubeObject(functionConfigMap)
// 		if functionConfigObject.GetKind() != "ConfigMap" {
// 			return nil, fmt.Errorf("input included functionConfig of unexpected kind %q; expected ConfigMap", functionConfigObject.GetKind())
// 		}
// 		data, found, err := functionConfigObject.GetMap("data")
// 		if err != nil {
// 			return nil, fmt.Errorf("error getting data: %w", err)
// 		}
// 		if found {
// 			entries, err := data.Entries()
// 			if err != nil {
// 				return nil, err
// 			}
// 			for k, v := range entries {
// 				sv, ok := fn.AsString(v)
// 				if !ok {
// 					return nil, fmt.Errorf("functionConfig key %q was not of type string", k)
// 				}
// 				functionConfig[k] = sv
// 			}
// 		}
// 	}

// 	out, err := fn(krmContext, items, functionConfig)
// 	if err != nil {
// 		return nil, fmt.Errorf("error from function: %w", err)
// 	}

// 	responseDoc := fn.NewDoc()
// 	response := responseDoc.NewMap()
// 	{
// 		response.SetString("apiVersion", "config.kubernetes.io/v1")
// 		response.SetString("kind", "ResourceList")

// 		items := response.UpsertSlice("items")

// 		for _, obj := range out {
// 			items.Add(obj.Node())
// 		}

// 		results := response.UpsertSlice("results")

// 		for _, result := range krmContext.Results {
// 			resultNode, err := fn.ObjectToMap(&result)
// 			if err != nil {
// 				return nil, fmt.Errorf("error converting result: %w", err)
// 			}
// 			results.Add(resultNode)
// 		}
// 	}

// 	responseBytes, err := responseDoc.ToYAML()
// 	if err != nil {
// 		return nil, fmt.Errorf("error building yaml response: %w", err)
// 	}

// 	return responseBytes, nil
// }

// type Result struct {
// 	Message  string   `json:"message,omitempty"`
// 	Severity Severity `json:"severity,omitempty"`
// 	// results:
// 	// - message: "Invalid type. Expected: integer, given: string"
// 	//   severity: error
// 	//   resourceRef:
// 	// 	apiVersion: v1
// 	// 	kind: Service
// 	// 	name: wordpress
// 	//   field:
// 	// 	path: spec.ports.0.port
// 	//   file:
// 	// 	path: service.yaml
// }

// // 	Severity is the severity of a result:
// type Severity string

// const (
// 	SeverityError   Severity = "error"
// 	SeverityWarning Severity = "warning"
// 	SeverityInfo    Severity = "info"
// )
