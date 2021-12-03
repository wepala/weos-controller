package weosgrpc_test

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

var port = ":8681"
var client pb.AccountClient

type GrpcTestAPI interface {
	weoscontroller.APIInterface
	HelloWorld(c context.Context) error
}

func setUpTest() (client pb.AccountClient, teardown func()) {

	InitalizeGrpc(ctx *context.Context, api weoscontroller.APIInterface,  "../fixtures/api/grpc.yaml")
	s := grpc.NewServer()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpc1 := &controller.GRPC{}
	err = grpc1.Initialize()
	if err != nil {
		log.Fatalf("Failed to intialize: %v", err)
	}
	pb.RegisterAccountServer(s, grpc1)
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
