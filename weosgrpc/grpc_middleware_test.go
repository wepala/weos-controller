//go:generate moq -out mocks_test.go -pkg weosgrpc_test . GrpcTestAPI
package weosgrpc_test

import (
	"context"
	"log"
	"net"

	weoscontroller "github.com/wepala/weos-controller"
	weosgrpc "github.com/wepala/weos-controller/weosgrpc"
	pb "github.com/wepala/weos-controller/weosgrpc/protofiles"
	"google.golang.org/grpc"
)

var port = ":8681"
var client pb.UserClient

type GrpcTestAPI interface {
	weoscontroller.GRPCAPIInterface
	SetAllMiddleware() error
	HelloWorld(c context.Context) error
	FooBar(c context.Context) error
}

func setUpTest() (client pb.UserClient, teardown func()) {

	//InitalizeGrpc(context.TODO(), api ,  "../fixtures/api/grpc.yaml")
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpc1 := &weosgrpc.GRPC{}
	err = grpc1.Initialize()
	if err != nil {
		log.Fatalf("Failed to intialize: %v", err)
	}
	pb.RegisterUserServer(s, grpc1)
	go func() {
		s.Serve(lis)
	}()

	conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client = pb.NewAccountClient(conn)
	teardown = func() {
		s.Stop()
		conn.Close()
		lis.Close()
	}

	return client, teardown
}
