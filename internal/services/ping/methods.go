package ping

import (
	"context"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"log/slog"
)

func (p *Ping) Ping(ctx context.Context, req *ssov1.PingRequest) (*ssov1.PingResponse, error) {
	const op = "ping.GetPing"

	log := p.log.With(slog.String("op", op))
	log.Info("Ping called with message", slog.String("msg", req.Message))

	return &ssov1.PingResponse{
		Reply: "Pong",
	}, nil
}
