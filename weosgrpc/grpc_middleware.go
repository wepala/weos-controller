package weosgrpc

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	weoscontroller "github.com/wepala/weos-controller"
	"google.golang.org/grpc"
)

type GrpcMiddleware struct {
	streamMiddleware grpc.ServerOption
	unaryMiddleware  grpc.ServerOption
}

func (g *GrpcMiddleware) setStreamMiddleware(args ...grpc.StreamServerInterceptor) {

	chainStream := grpc_middleware.ChainStreamServer(args...)
	g.streamMiddleware = grpc.StreamInterceptor(chainStream)
}

func (g *GrpcMiddleware) setUnaryMiddleware(args ...grpc.UnaryServerInterceptor) {

	chainUnary := grpc_middleware.ChainUnaryServer(args...)
	g.streamMiddleware = grpc.UnaryInterceptor(chainUnary)
}

func (g *GrpcMiddleware) GetStreamMiddleware() grpc.ServerOption {
	return g.streamMiddleware
}

func (g *GrpcMiddleware) GetUnaryMiddleware() grpc.ServerOption {
	return g.unaryMiddleware
}

func SetAllMiddleware(ctx *context.Context, config *weoscontroller.APIConfig) *context.Context {
	grpcStream := make([]grpc.StreamServerInterceptor, 2)
	grpcUnary := make([]grpc.UnaryServerInterceptor, 2)
	//TODO call the functions to convert the middleware to the interceptors and append to array
	//call setUnaryMiddleware and setStreamMiddleware with the array

	grpcMiddlewareConfig := config.Grpc.Middlewares

	for _, streamMiddleware := range grpcMiddlewareConfig.Stream.Middleware {
		switch streamMiddleware {
		case "Authenticate":
			grpcStream = append(grpcStream, grpc_auth.StreamServerInterceptor(Authenticate))
		case "Recovery":
			grpcStream = append(grpcStream, grpc_recovery.StreamServerInterceptor())
		}
	}

	for _, UnaryMiddleware := range grpcMiddlewareConfig.Unary.Middleware {
		switch UnaryMiddleware {
		case "Authenticate":
			grpcUnary = append(grpcUnary, grpc_auth.UnaryServerInterceptor(Authenticate))
		case "Recovery":
			grpcUnary = append(grpcUnary, grpc_recovery.UnaryServerInterceptor())
		}
	}

	var grpcMiddleware GrpcMiddleware

	grpcMiddleware.setStreamMiddleware(grpcStream...)
	grpcMiddleware.setUnaryMiddleware(grpcUnary...)

	//WithValue it into the context? not sure
	context := context.WithValue(*ctx, "grpcServerOptions", grpcMiddleware)
	return &context
}
