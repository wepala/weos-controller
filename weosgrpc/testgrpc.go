package weosgrpc

import (
	"net/http"
	"time"

	pb "github.com/wepala/weoscontroller/weosgrpc/protofiles"
	"golang.org/x/net/context"
)

type GRPC struct {
	pb.UnimplementedAccountServer
	Client *http.Client
}

func (g *GRPC) CreateUser(ctxt context.Context, a *pb.Request) (*pb.Response, error) {
	accountPayload := model.AccountPayload{
		ID:    a.ID,
		Title: a.Title,
	}
	// Dispatch command to create accountPayload with payload
	if accountPayload.ID == "" {
		accountPayload.ID = GenerateID()
	}

	return &pb.Response{Account: a, Result: "account created successfully", IsValid: true}, nil
}

//Initialize and setup configurations for GRPC controller
func (g *GRPC) InitializeTestGrpc() error {
	var err error
	//initialize app
	if g.Client == nil {
		g.Client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	return nil
}
