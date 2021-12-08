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

func setUpTest(api *GrpcTestAPIMock) (client pb.UserClient, teardown func(), ctx context.Context) {

	weosgrpc.InitalizeGrpc(ctx, api, "../fixtures/api/grpc.yaml")
	s := grpc.NewServer(api.GetStreamMiddleware, api.GetUnaryMiddleware)
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

	return client, teardown, ctx
}

func TestGrpcMiddleware(t *testing.T) {

	var middlewareAndHandlersCalled []string
	api := &GrpcTestAPIMock{
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		ContextFunc: func() context.Context {
			return nil
		},
		SetContextFunc: func(c context.Context) {},
		FooBarFunc: func(c context.Context) error {
			middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "fooBarHandler")
			return nil
		},
		HelloWorldFunc: func(c context.Context) error {
			middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "helloWorldHandler")
			return nil
		},
		InitializeFunc: func() error {
			return nil
		},
		SetAllMiddlewareFunc: func() {},
		GetStreamMiddlewareFunc: func() grpc.ServerOption {
			middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "getStreamMiddleware")
			return nil
		},
		GetUnaryMiddlewareFunc: func() grpc.ServerOption {
			middlewareAndHandlersCalled = append(middlewareAndHandlersCalled, "getUnaryMiddleware")
			return nil
		},
	}
	client, teardown, ctx := setUpTest(api)

	defer teardown()
	req := &pb.Request{
		ID:    "123sa",
		Title: "Account132",
	}
	_, err := client.CreateUser(ctx, req)
	if err != nil {
		t.Errorf("unexpected error creating user  %s", err)
	}
	//check that the expected handlers and middleware are called
	if len(middlewareAndHandlersCalled) != 4 {
		t.Fatalf("expected %d middlewares and handlers to be called, got %d", 4, len(middlewareAndHandlersCalled))
	}

}

func TestCreateUserGRPC(t *testing.T) {

	ctx := context.Background()

	api := &GrpcTestAPIMock{
		AddConfigFunc: func(config *weoscontroller.APIConfig) error {
			return nil
		},
		AddPathConfigFunc: func(path string, config *weoscontroller.PathConfig) error {
			return nil
		},
		ContextFunc: func() context.Context {
			return ctx
		},
		SetContextFunc: func(c context.Context) {
			ctx = c
		},
		FooBarFunc: func(c context.Context) error {
			return nil
		},
		HelloWorldFunc: func(c context.Context) error {
			return nil
		},
		InitializeFunc: func() error {
			return nil
		},
		SetAllMiddlewareFunc: func() {},
		GetStreamMiddlewareFunc: func() grpc.ServerOption {
			return nil
		},
		GetUnaryMiddlewareFunc: func() grpc.ServerOption {
			return nil
		},
	}

	client, teardown, ctx := setUpTest(api)

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
		if resp.Result != "user created successfully" {
			t.Errorf("unexpected error, expected result to %s got %s", "user created successfully", resp.Result)
		}
		if resp.User.ID != req.ID {
			t.Errorf("unexpected error, expected user id to %s got %s", req.ID, resp.User.ID)
		}
		if resp.User.Title != req.Title {
			t.Errorf("unexpected error, expected user title to %s got %s", req.Title, resp.User.Title)
		}
	})
}
