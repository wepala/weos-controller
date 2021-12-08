package weosgrpc

import (
	"context"
	"net/http"
	"time"

	pb "github.com/wepala/weos-controller/weosgrpc/protofiles"
)

type Grpc struct {
	pb.UnimplementedUserServer
	Client *http.Client
}

func (g *Grpc) CreateUser(ctxt context.Context, a *pb.Request) (*pb.Response, error) {
	return &pb.Response{User: a, Result: "account created successfully", IsValid: true}, nil
}

//Initialize and setup configurations for GRPC controller
func (g *Grpc) Initialize() error {
	//initialize app
	if g.Client == nil {
		g.Client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
	return nil
}
