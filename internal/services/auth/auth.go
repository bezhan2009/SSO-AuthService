package auth

import (
	"SSO/internal/domain/models"
	"SSO/internal/lib/jwt"
	"SSO/internal/lib/logger/sl"
	"SSO/internal/storage"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidAppId       = errors.New("invalid app id")
	ErrUserExists         = errors.New("user already exists")
	ErrEmailIsRequired    = errors.New("email is required")
	ErrNameIsRequired     = errors.New("name is required")
	ErrPasswordIsRequired = errors.New("password is required")
)

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context,
		email string,
		passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context,
		email string) (models.User, error)
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
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:         log,
		usrSaver:    usrSaver,
		usrProvider: usrProvider,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context,
	email string,
	password string,
	appID int) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Logging user")

	if err = validateUserRequest(email, password); err != nil {
		log.Warn("Validation Error", slog.String("error", err.Error()))

		return
	}

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found")

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		a.log.Error("Failed to get user", "error", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Warn("Invalid credentials", sl.Err(err))

		return "", ErrInvalidCredentials
	}

	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		a.log.Error("Failed to get app", "error", err)

		return "", ErrInvalidCredentials
	}

	log.Info("User logged in successfully")

	token, err = jwtauth.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("Failed to create token", "error", err)

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context,
	email string,
	password string) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("Registering user")

	if err = validateUserRequest(email, password); err != nil {
		log.Warn("Validation Error", slog.String("error", err.Error()))

		return
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to GenerateFromPassword", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uid, err := a.usrSaver.SaveUser(ctx, email, passHash)
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

func validateUserRequest(email, password string) (err error) {
	if email == "" {
		return ErrEmailIsRequired
	}

	if password == "" {
		return ErrPasswordIsRequired
	}

	return nil
}
