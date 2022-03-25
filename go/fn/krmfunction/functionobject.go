package krmfunction

import (
	"context"
	"fmt"
	"reflect"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
)

const (
	ExportedConfigMapName = "exports.kptfile.kpt.dev"
	ImportedConfigMapName = "kptfile.kpt.dev"
)

type KRMFunctionObject interface {
	Run(ctx *fn.Context, in []*fn.KubeObject, config map[string]string) ([]*fn.KubeObject, error)
}

type processor struct {
	fnObject KRMFunctionObject
}

func ConvertObjectToResourceListProcessor(fnObject KRMFunctionObject) fn.ResourceListProcessor {
	if fnObject == nil {
		panic("fnObject was nil")
	}
	return &processor{fnObject: fnObject}
}

func (p *processor) Process(rl *fn.ResourceList) error {
	krmContext := &fn.Context{Context: context.Background()}

	functionConfig := map[string]string{}
	if rl.FunctionConfig != nil {
		if rl.FunctionConfig.GetKind() != "ConfigMap" {
			return fmt.Errorf("input included functionConfig of unexpected kind %q; expected ConfigMap", rl.FunctionConfig.GetKind())
		}
		data := rl.FunctionConfig.UpsertMap("data")

		entries, err := data.StringStringEntries()
		if err != nil {
			return err
		}
		for k, v := range entries {
			functionConfig[k] = v
		}
	}

	out, err := p.run(krmContext, rl.Items, functionConfig)
	if err != nil {
		return err
	}
	rl.Items = out

	rl.Results = append(rl.Results, krmContext.Results...)

	return nil
}

func (p *processor) run(ctx *fn.Context, in []*fn.KubeObject, fnConfig map[string]string) ([]*fn.KubeObject, error) {
	fnValue := reflect.ValueOf(p.fnObject).Elem()
	fnType := fnValue.Type()

	var imported *fn.KubeObject
	for _, obj := range in {
		apiVersion := obj.GetAPIVersion()
		kind := obj.GetKind()
		if apiVersion == "v1" && kind == "ConfigMap" {
			if obj.GetName() == ImportedConfigMapName {
				imported = obj
			}
		}
	}

	importedData := make(map[string]string)
	if imported != nil {
		// TODO: support non-ConfigMap?
		entries, err := imported.UpsertMap("data").StringStringEntries()
		if err != nil {
			return nil, fmt.Errorf("error parsing configmap data: %w", err)
		}
		for key, v := range entries {
			importedData[key] = v
		}
	}

	fqn := importedData["fqn"]
	if fqn == "" {
		// TODO: Move to porch (make this a first class concept?)
		parentFQN := importedData["parent.fqn"]
		if parentFQN == "" {
			if len(importedData) == 1 {
				// HACK: Assuming this is a root package with no parent
				parentFQN = ""
			} else {
				return nil, fmt.Errorf("parent.fqn not found in importedData (for computing fqn)")
			}
		}
		name := importedData["name"]
		if name == "" {
			return nil, fmt.Errorf("name not found in importedData (for computing fqn)")
		}
		fqn = parentFQN
		if fqn != "" {
			fqn += "/"
		}
		fqn += name
	}

	var exports reflect.Value
	for i := 0; i < fnType.NumField(); i++ {
		field := fnType.Field(i)
		bind := field.Tag.Get("bind")
		if bind == "" {
			if !field.IsExported() {
				continue
			}

			switch field.Name {
			case "Name":
				name := importedData["name"]
				if name == "" {
					return nil, fmt.Errorf("name not found in importedData")
				}
				fnValue.Field(i).SetString(name)
			case "FQN":
				fnValue.Field(i).SetString(fqn)
			case "Imports":
				if err := bindImports(fnValue.Field(i), importedData); err != nil {
					return nil, fmt.Errorf("error binding imports: %w", err)
				}
			case "Exports":
				exports = fnValue.Field(i)
				// We map these "out" later
			default:
				return nil, fmt.Errorf("unknown field %q", field.Name)
			}
		}
	}

	out, err := p.fnObject.Run(ctx, in, fnConfig)
	if err != nil {
		return nil, err
	}

	shouldExport := true // We always export the fqn
	if shouldExport {
		var exported *fn.KubeObject
		for _, obj := range in {
			apiVersion := obj.GetAPIVersion()
			kind := obj.GetKind()
			if apiVersion == "v1" && kind == "ConfigMap" {
				if obj.GetName() == ExportedConfigMapName {
					exported = obj
				}
			}
		}

		if exported == nil {
			exported = fn.NewKubeObject()
			exported.SetName(ExportedConfigMapName)
			exported.SetAPIVersion("v1")
			exported.SetKind("ConfigMap")
			exported.SetAnnotation("config.kubernetes.io/local-config", "true")
			out = append(out, exported)
		}

		exportsMap := exported.UpsertMap("data")
		if exports.IsValid() {
			if err := bindExports(exports, exportsMap); err != nil {
				return nil, fmt.Errorf("error binding exports: %w", err)
			}
		}
		exportsMap.SetString("fqn", fqn)
	}

	return out, err
}

func bindImports(importsValue reflect.Value, importedData map[string]string) error {
	importsType := importsValue.Type()
	for i := 0; i < importsType.NumField(); i++ {
		field := importsType.Field(i)
		bind := field.Tag.Get("bind")
		if bind == "" {
			switch field.Name {
			case "Folder":
				folderName := importedData["parent.folder.name"]
				ref := &ConcreteObjectReference{
					Name: folderName,
				}
				importsValue.Field(i).Set(reflect.ValueOf(ref))
			case "Namespace":
				namespaceName := importedData["parent.namespace.name"]
				ref := &ConcreteObjectReference{
					Name: namespaceName,
				}
				importsValue.Field(i).Set(reflect.ValueOf(ref))
			default:
				return fmt.Errorf("unknown imports field %q", field.Name)
			}
		}
	}
	return nil
}

type ObjectRef interface {
	GetName() string
}

func bindExports(exportsValue reflect.Value, data *fn.SubObject) error {
	exportsType := exportsValue.Type()
	for i := 0; i < exportsType.NumField(); i++ {
		field := exportsType.Field(i)
		val := exportsValue.Field(i).Interface()

		if val == nil {
			continue
		}

		bind := field.Tag.Get("bind")
		if bind == "" {
			switch field.Name {
			case "Folder":
				ref, ok := val.(ObjectRef)
				if !ok {
					return fmt.Errorf("expected ObjectRef, was %T", val)
				}
				data.SetString("folder.name", ref.GetName())
			case "Namespace":
				ref, ok := val.(ObjectRef)
				if !ok {
					return fmt.Errorf("expected ObjectRef, was %T", val)
				}
				data.SetString("namespace.name", ref.GetName())
			default:
				return fmt.Errorf("unknown exports field %q", field.Name)
			}
		}
	}
	return nil
}

type ConcreteObjectReference struct {
	Name string
}

func (r *ConcreteObjectReference) GetName() string {
	return r.Name
}
