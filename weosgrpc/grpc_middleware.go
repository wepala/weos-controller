package weosgrpc

import (
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

	/*g.streamMiddleware = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		array...

		//grpc_zap.StreamServerInterceptor(zapLogger),
		grpc_auth.StreamServerInterceptor(myAuthFunction), // myAuthFunction = Authenticate function in weosgrpc_api.go
		grpc_recovery.StreamServerInterceptor(),
	))*/
}

func (g *GrpcMiddleware) setUnaryMiddleware(args ...grpc.UnaryServerInterceptor) {

	/*zapLogger := &zap.Logger{}
	g.unaryMiddleware = grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		array...
		grpc_zap.UnaryServerInterceptor(zapLogger),
		grpc_auth.UnaryServerInterceptor(myAuthFunction), // myAuthFunction = Authenticate function in weosgrpc_api.go
		grpc_recovery.UnaryServerInterceptor(),
	))*/
}

func (g *GrpcMiddleware) GetStreamMiddleware() grpc.ServerOption {
	return g.streamMiddleware
}

func (g *GrpcMiddleware) GetUnaryMiddleware() grpc.ServerOption {
	return g.unaryMiddleware
}

func SetAllMiddleware(config *weoscontroller.APIConfig) {
	grpcStream := make([]grpc.StreamServerInterceptor, 2)
	grpcUnary := make([]grpc.UnaryServerInterceptor, 2)
	//TODO call the functions to convert the middleware to the interceptors and append to array
	//call setUnaryMiddleware and setStreamMiddleware with the array

	grpcMiddleware := config.Grpc.Middlewares

	for _, streamMiddleware := range grpcMiddleware.Stream.Middleware {
		switch streamMiddleware {
		case "Authenticate":
			grpcStream = append(grpcStream, grpc_auth.StreamServerInterceptor(Authenticate))
		case "Recovery":
			grpcStream = append(grpcStream, grpc_recovery.StreamServerInterceptor())
		}
	}

	for _, UnaryMiddleware := range grpcMiddleware.Unary.Middleware {
		switch UnaryMiddleware {
		case "Authenticate":
			grpcUnary = append(grpcUnary, grpc_auth.UnaryServerInterceptor(Authenticate))
		case "Recovery":
			grpcUnary = append(grpcUnary, grpc_recovery.UnaryServerInterceptor())
		}
	}

	//Setup Errors
}
