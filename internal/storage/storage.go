package storage

import (
	"SSO/internal/domain/models"
	"context"
	"errors"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
)

type Storage interface {
	SaveUser(ctx context.Context, userRequest *ssov1.RegisterRequest) (int64, error)
	User(ctx context.Context, username string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
	App(ctx context.Context, appID int64) (models.App, error)
}

var (
	ErrUserExists    = errors.New("user already exists")
	ErrUserNotFound  = errors.New("user not found")
	ErrAdminNotFound = errors.New("admin not found")
	ErrAppNotFound   = errors.New("app not found")
)
