//go:generate moq -out grpc_mocks_test.go -pkg weosgrpc_test . GrpcTestAPI

package weoscontroller

import "golang.org/x/net/context"

type GrpcTestAPI interface {
	GRPCAPIInterface
	HelloWorld(c context.Context) error
	FooBar(c context.Context) error
}
