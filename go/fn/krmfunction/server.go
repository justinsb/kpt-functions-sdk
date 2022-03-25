package krmfunction

// import (
// 	"context"
// 	"fmt"
// 	"net"
// 	"strings"

// 	pb "github.com/GoogleContainerTools/kpt/porch/func/evaluator"
// 	"google.golang.org/grpc"
// 	"k8s.io/klog/v2"
// )

// type functionEvaluator struct {
// 	registry *Registry

// 	pb.UnimplementedFunctionEvaluatorServer
// }

// func (r *Registry) runServer(listen string) error {
// 	lis, err := net.Listen("tcp", listen)
// 	if err != nil {
// 		return fmt.Errorf("failed to listen on %q: %w", listen, err)
// 	}

// 	evaluator := &functionEvaluator{
// 		registry: r,
// 	}

// 	klog.Infof("Listening on %s", listen)

// 	// Start the gRPC server
// 	server := grpc.NewServer()
// 	pb.RegisterFunctionEvaluatorServer(server, evaluator)
// 	if err := server.Serve(lis); err != nil {
// 		return fmt.Errorf("server failed: %w", err)
// 	}
// 	return nil
// }

// func lastComponent(s string) string {
// 	lastSlash := strings.LastIndex(s, "/")
// 	return s[lastSlash+1:]
// }

// func (e *functionEvaluator) EvaluateFunction(ctx context.Context, req *pb.EvaluateFunctionRequest) (*pb.EvaluateFunctionResponse, error) {
// 	image := req.Image

// 	tokens := strings.SplitN(image, ":", 2)
// 	if len(tokens) != 2 {
// 		// TODO: Assume latest?
// 		return nil, fmt.Errorf("expected version in image %q", image)
// 	}

// 	functionName := lastComponent(tokens[0])

// 	krmFunctionObject := e.registry.FindFunc(functionName)
// 	if krmFunctionObject == nil {
// 		return nil, fmt.Errorf("function %q not found (for image %q)", functionName, req.Image)
// 	}

// 	krmFunction := ConvertObjectToFunction(krmFunctionObject)
// 	krmContext := &Context{Context: ctx}
// 	stdout, err := Run(krmContext, krmFunction, req.ResourceList)
// 	if err != nil {
// 		return nil, fmt.Errorf("error running function %q: %w", functionName, err)
// 	}

// 	klog.Infof("Evaluated %q: stdout %d bytes", req.Image, len(stdout))

// 	// TODO: include stderr in the output?
// 	return &pb.EvaluateFunctionResponse{
// 		ResourceList: stdout,
// 	}, nil
// }
