//go:generate moq -out mocks_test.go -pkg weosgrpc_test . GrpcTestAPI
//cannot generate in folder, need be generated outside and then moved in
package weosgrpc_test

import (
	"context"
	"log"
	"net"
	"testing"

	weoscontroller "github.com/wepala/weos-controller"
	weosgrpc "github.com/wepala/weos-controller/weosgrpc"
	pb "github.com/wepala/weos-controller/weosgrpc/protofiles"
	"google.golang.org/grpc"
)

var port = ":8681"
var client pb.UserClient

type GrpcTestAPI interface {
	weoscontroller.GRPCAPIInterface
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
	grpc1 := &weosgrpc.Grpc{}
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

	client = pb.NewUserClient(conn)
	teardown = func() {
		s.Stop()
		conn.Close()
		lis.Close()
	}

	return client, teardown
}

func TestCreateAccountGRPC(t *testing.T) {

	//initialization will instantiate with application so we need to overwrite with our mock application
	client, teardown := setUpTest()

	ctx := context.Background()
	defer teardown()
	t.Run("basic request", func(t *testing.T) {
		req := &pb.Request{
			ID:    "123sa",
			Title: "Account132",
		}
		resp, err := client.CreateUser(ctx, req)
		if err != nil {
			t.Errorf("unexpected error creating user  %s", err)
		}
		if resp.IsValid != true {
			t.Errorf("unexpected error, expected isvalid to be true but got false")
		}
		if resp.Result != "account created successfully" {
			t.Errorf("unexpected error, expected result to %s got %s", "account created successfully", resp.Result)
		}
		if resp.User.ID != req.ID {
			t.Errorf("unexpected error, expected user id to %s got %s", req.ID, resp.User.ID)
		}
		if resp.User.Title != req.Title {
			t.Errorf("unexpected error, expected user title to %s got %s", req.Title, resp.User.Title)
		}
	})
}
