package auth

import (
	"SSO/internal/domain/models"
	"context"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"google.golang.org/grpc"
)

type Auth interface {
	Login(ctx context.Context,
		userLogin *ssov1.LoginRequest) (tokenResponse models.TokenResponse, err error)
	RegisterNewUser(ctx context.Context,
		user *ssov1.RegisterRequest) (userID int64, err error)
	IsAdmin(ctx context.Context,
		userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(grpc *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(grpc, &serverAPI{auth: auth})
}
