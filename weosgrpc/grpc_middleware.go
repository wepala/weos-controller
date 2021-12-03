package weosgrpc

import (
	weoscontroller "github.com/wepala/weos-controller"
	"google.golang.org/grpc"
)

type GrpcMiddleware struct {
	streamMiddleware grpc.ServerOption
	unaryMiddleware  grpc.ServerOption
}

func (g *GrpcMiddleware) setStreamMiddleware(array *[]grpc.StreamServerInterceptor) {

	/*g.streamMiddleware = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		array...

		//grpc_zap.StreamServerInterceptor(zapLogger),
		grpc_auth.StreamServerInterceptor(myAuthFunction), // myAuthFunction = Authenticate function in weosgrpc_api.go
		grpc_recovery.StreamServerInterceptor(),
	))*/
}

func (g *GrpcMiddleware) setUnaryMiddleware(array *[]grpc.UnaryServerInterceptor) {

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
	var grpcSteam *[]grpc.StreamServerInterceptor
	var grpcUnary *[]grpc.UnaryServerInterceptor
	//TODO call the functions to convert the middleware to the interceptors and append to array
	//call setUnaryMiddleware and setStreamMiddleware with the array

	grpcMiddleware := config.Grpc.Middlewares

	for _, streamMiddleware := range grpcMiddleware.Stream.Middleware {
		switch streamMiddleware {
		case "Authenticate":

		case "Recovery":
		}
	}

}
