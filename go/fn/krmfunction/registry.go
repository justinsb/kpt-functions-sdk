package krmfunction

import (
	"flag"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn"
	"k8s.io/klog/v2"
)

// import (
// 	"flag"
// 	"fmt"
// 	"os"
// 	"reflect"
// 	"strings"
// )

type Registry struct {
	functions map[string]KRMFunctionObject
}

func (r *Registry) Register(fn KRMFunctionObject) {
	fnValue := reflect.ValueOf(fn).Elem()
	name := fnValue.Type().Name()
	r.registerWithName(name, fn)
}

func (r *Registry) registerWithName(name string, fn KRMFunctionObject) {
	name = r.normalizeName(name)
	if r.functions == nil {
		r.functions = make(map[string]KRMFunctionObject)
	}
	r.functions[name] = fn
}

func (r *Registry) normalizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	return name
}

func (r *Registry) FindFunc(name string) KRMFunctionObject {
	name = r.normalizeName(name)
	fn := r.functions[name]
	return fn
}

func (r *Registry) RunAsMain() {
	klog.InitFlags(nil)
	// listen := ""
	// flag.StringVar(&listen, "listen", listen, "run in server mode, listening on the given endpoint")
	flag.Parse()

	// if listen != "" {
	// 	if err := r.runServer(listen); err != nil {
	// 		fmt.Fprintf(os.Stderr, "error running server: %v\n", err)
	// 		os.Exit(1)
	// 	}
	// 	return
	// }

	arg0 := os.Args[0]
	fnObject := r.FindFunc(arg0)
	if fnObject == nil {
		var knownFunctions []string
		for k := range r.functions {
			knownFunctions = append(knownFunctions, k)
		}
		fmt.Fprintf(os.Stderr, "unable to find function %q (known %s)\n", arg0, strings.Join(knownFunctions, ","))
		os.Exit(1)
	}

	p := ConvertObjectToResourceListProcessor(fnObject)
	if err := fn.AsMain(p); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func RunAsMain(fnObject KRMFunctionObject) {
	r := &Registry{}
	name := os.Args[0]
	name = r.normalizeName(name)

	r.registerWithName(name, fnObject)
	r.RunAsMain()
}

// type functionWrapper struct {
// 	function func(ctx *fn.Context, in []*fn.KubeObject, config map[string]string) ([]*fn.KubeObject, error)
// }

// func (w *functionWrapper) Run(ctx *fn.Context, in []*fn.KubeObject, config map[string]string) ([]*fn.KubeObject, error) {
// 	return w.function(ctx, in, config)
// }

// func Wrap(function func(ctx *fn.Context, in []*fn.KubeObject, config map[string]string) ([]*fn.KubeObject, error)) KRMFunctionObject {
// 	return &functionWrapper{function: function}
// }
