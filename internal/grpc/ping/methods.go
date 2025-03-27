package ping

import (
	"context"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
)

func (s *ServerAPI) Ping(ctx context.Context, req *ssov1.PingRequest) (*ssov1.PingResponse, error) {
	return s.ping.Ping(ctx, req)
}
