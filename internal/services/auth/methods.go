package auth

import (
	"SSO/internal/domain/models"
	jwtauth "SSO/internal/lib/jwt"
	"SSO/internal/lib/logger/sl"
	"SSO/internal/services/auth/validators"
	"SSO/internal/storage"
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func (a *Auth) Login(ctx context.Context,
	userLogin *ssov1.LoginRequest) (tokenResponse models.TokenResponse, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", userLogin.GetUsername()),
	)

	log.Info("Logging user")

	if err = validators.ValidateUserLoginRequest(userLogin); err != nil {
		log.Warn("Validation Error", slog.String("error", err.Error()))

		return tokenResponse, err
	}

	user, err := a.usrProvider.User(ctx, userLogin.GetUsername())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found")

			return tokenResponse, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("Failed to get user", "error", err)
		return tokenResponse, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(userLogin.GetPassword())); err != nil {
		a.log.Warn("Invalid credentials", sl.Err(err))

		return tokenResponse, ErrInvalidCredentials
	}

	app, err := a.appProvider.App(ctx, int64(userLogin.GetAppLogin()))
	if err != nil {
		a.log.Error("Failed to get app", slog.String("error", err.Error()))

		return tokenResponse, ErrInvalidCredentials
	}

	log.Info("User logged in successfully")

	accessToken, refreshToken, err := jwtauth.NewToken(user,
		app,
		a.authConfig.JwtTTLMinutes,
		a.authConfig.JwtTTLRefreshHours,
	)
	if err != nil {
		a.log.Error("Failed to create token", slog.String("error", err.Error()))

		return tokenResponse, fmt.Errorf("%s: %w", op, err)
	}

	tokenResponse.AccessToken = accessToken
	tokenResponse.RefreshToken = refreshToken
	tokenResponse.UserID = int(user.ID)

	return tokenResponse, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context,
	user *ssov1.RegisterRequest) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", user.GetEmail()),
	)

	log.Info("Registering user")

	if err = validators.ValidateUserRegisterRequest(user); err != nil {
		log.Warn("Validation Error", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to GenerateFromPassword", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	user.Password = string(passHash)

	uid, err := a.usrSaver.SaveUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			a.log.Warn("User already exists")

			return 0, ErrUserExists
		}

		log.Error("Failed to SaveUser", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User Registered")

	return uid, nil
}

func (a *Auth) IsAdmin(ctx context.Context,
	userID int64) (isAdmin bool, err error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", int(userID)),
	)

	log.Info("Checking if user is admin")

	isAdmin, err = a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("App not found")

			return false, ErrInvalidAppId
		}

		a.log.Error("Failed to check if user is admin", "error", err)

		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
