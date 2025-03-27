package ping

import (
	"context"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"google.golang.org/grpc"
)

type PingAPI interface {
	Ping(ctx context.Context, req *ssov1.PingRequest) (*ssov1.PingResponse, error)
}

type ServerAPI struct {
	ssov1.UnimplementedPingServiceServer
	ping PingAPI
}

func Register(grpc *grpc.Server, ping PingAPI) {
	ssov1.RegisterPingServiceServer(grpc, &ServerAPI{ping: ping})
}
