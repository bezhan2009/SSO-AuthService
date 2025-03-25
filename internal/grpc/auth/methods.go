package auth

import (
	"SSO/internal/grpc/handlers"
	"SSO/internal/services/auth"
	"SSO/internal/storage"
	"context"
	"errors"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	tokenResponse, err := s.auth.Login(ctx, req)
	if err != nil {
		return nil, handlers.HandleError(err)
	}

	return &ssov1.LoginResponse{
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		UserId:       int64(tokenResponse.UserID),
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, handlers.HandleError(err)
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	if err := validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, int64(req.GetUserId()))
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "app not found")
		}

		return nil, handlers.HandleError(err)
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
