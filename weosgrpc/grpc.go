package weosgrpc

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	pb "github.com/wepala/weos-controller/weosgrpc/protofiles"
)

type Grpc struct {
	pb.UnimplementedUserServer
	Client *http.Client
	DB     *sql.DB
}

func (g *GRPC) CreateUser(ctxt context.Context, a *pb.Request) (*pb.Response, error) {
	return &pb.Response{User: a, Result: "account created successfully", IsValid: true}, nil
}

func (g *GRPC) Initialize() {
	if g.Client == nil {
		g.Client = &http.Client{
			Timeout: time.Second * 10,
		}
	}
}
