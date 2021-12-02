package weosgrpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/uber-go/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GrpcMiddleware struct {
	streamMiddleware grpc.ServerOption
	unaryMiddleware  grpc.ServerOption
}

func (g *GrpcMiddleware) SetStreamMiddleware(zapLogger *zap.Logger) {

	g.streamMiddleware = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		//grpc_zap.StreamServerInterceptor(zapLogger),
		grpc_auth.StreamServerInterceptor(myAuthFunction), // myAuthFunction = Authenticate function in weosgrpc_api.go
		grpc_recovery.StreamServerInterceptor(),
	))
}

func (g *GrpcMiddleware) SetUnaryMiddleware() {

	zapLogger := &zap.Logger{}
	g.unaryMiddleware = grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		grpc_zap.UnaryServerInterceptor(zapLogger),
		grpc_auth.UnaryServerInterceptor(myAuthFunction), // myAuthFunction = Authenticate function in weosgrpc_api.go
		grpc_recovery.UnaryServerInterceptor(),
	))
}

func (g *GrpcMiddleware) GetStreamMiddleware() grpc.ServerOption {
	return g.streamMiddleware
}

func (g *GrpcMiddleware) GetUnaryMiddleware() grpc.ServerOption {
	return g.unaryMiddleware
}
