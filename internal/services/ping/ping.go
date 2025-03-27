package ping

import (
	"context"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"log/slog"
)

type Ping struct {
	log *slog.Logger
}

type PingAPI interface {
	Ping(ctx context.Context, req *ssov1.PingRequest) (*ssov1.PingResponse, error)
}

func New(log *slog.Logger) *Ping {
	return &Ping{
		log: log,
	}
}
