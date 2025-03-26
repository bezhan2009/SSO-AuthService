package postgres

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage"
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"github.com/jackc/pgx/v5/pgconn"
	"log/slog"

	"gorm.io/gorm"
)

type Storage struct {
	db  *gorm.DB
	log *slog.Logger
}

// New создает новый экземпляр хранилища с PostgreSQL.
func New(db *gorm.DB, log *slog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	// Автоматическая миграция моделей
	if err := db.AutoMigrate(&models.User{}, &models.App{}, &models.Admin{}); err != nil {
		log.Error(fmt.Sprintf("Error migrating database: %s", err), slog.String("error", err.Error()))
		return nil, err
	}

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) SaveUser(ctx context.Context, userRequest *ssov1.RegisterRequest) (models.User, error) {
	const op = "storage.postgres.SaveUser"

	user := models.User{
		FirstName:    userRequest.GetFirstName(),
		LastName:     userRequest.GetLastName(),
		Email:        userRequest.GetEmail(),
		Username:     userRequest.GetUsername(),
		HashPassword: userRequest.GetPassword(),
	}
	result := s.db.WithContext(ctx).Create(&user)

	if result.Error != nil {
		// Проверяем, является ли ошибка ошибкой уникального ограничения
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			s.log.Warn("User already exists", slog.String("err", result.Error.Error()))
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		s.log.Error("Error inserting user", slog.String("err", result.Error.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

func (s *Storage) User(ctx context.Context, username string) (models.User, error) {
	const op = "storage.postgres.User"

	var user models.User
	result := s.db.WithContext(ctx).Where("email = ? OR username = ?", username, username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			s.log.Warn("User not found", slog.String("username", username))
			return models.User{}, storage.ErrUserNotFound
		}
		s.log.Error("Error fetching user", slog.String("err", result.Error.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var admin models.Admin
	result := s.db.WithContext(ctx).Select("id").Where("user_id = ?", userID).First(&admin)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAdminNotFound)
		}
		s.log.Error("Error fetching admin", slog.String("err", result.Error.Error()))
		return false, fmt.Errorf("%s: %w", op, result.Error)
	}

	return true, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.postgres.App"

	var app models.App
	result := s.db.WithContext(ctx).Where("id = ?", appID).First(&app)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.App{}, storage.ErrAppNotFound
		}

		s.log.Error("Error fetching app", slog.String("err", result.Error.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return app, nil
}

func (s *Storage) Apps() {

}
