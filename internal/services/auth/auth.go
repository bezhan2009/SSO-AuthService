package auth

import (
	"SSO/internal/config"
	"SSO/internal/domain/models"
	"context"
	"errors"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"log/slog"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	authConfig  config.AuthParams
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		userRequest *ssov1.RegisterRequest) (user models.User, err error)
}

type UserProvider interface {
	User(ctx context.Context,
		username string) (models.User, error)
	IsAdmin(ctx context.Context,
		userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context,
		appID int64) (models.App, error)
}

// New returns a new instance of the Auth service.
func New(log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	appProvider AppProvider,
	authConfig config.AuthParams) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvider: appProvider,
		authConfig:  authConfig,
	}
}
