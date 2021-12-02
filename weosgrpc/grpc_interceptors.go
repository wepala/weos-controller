package weosgrpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
)

type GrpcStreamMiddleware struct {
	streamMiddleware grpc.ServerOption
	unaryMiddleware  grpc.ServerOption
}

func (g *GrpcStreamMiddleware) SetStreamMiddleware() error {
	if g.streamMiddleware == nil {
		g.streamMiddleware = grpc.ServerOption{}
	}
	g.streamMiddleware = grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_opentracing.StreamServerInterceptor(),
		grpc_prometheus.StreamServerInterceptor,
		grpc_zap.StreamServerInterceptor(zapLogger),
		grpc_auth.StreamServerInterceptor(myAuthFunction),
		grpc_recovery.StreamServerInterceptor(),
	))
	return nil
}
