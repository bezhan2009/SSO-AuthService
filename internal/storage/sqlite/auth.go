package sqlite

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage"
	"context"
	"errors"
	"fmt"
	ssov1 "github.com/bezhan2009/AuthProtos/gen/go/sso"
	"github.com/mattn/go-sqlite3"
	"log/slog"

	"gorm.io/gorm"
)

type Storage struct {
	db  *gorm.DB
	log *slog.Logger
}

// New создает новый экземпляр хранилища с SQLite.
func New(db *gorm.DB, log *slog.Logger) (*Storage, error) {
	const op = "storage.sqlite.New"

	// Автоматическая миграция моделей
	if err := db.AutoMigrate(&models.User{}, &models.App{}, &models.Admin{}); err != nil {
		log.Error(fmt.Sprintf("op: %s: Error migrating database", op), slog.String("error", err.Error()))
		return nil, err
	}

	if db.Where("id = 1").First(&models.App{}).Error != nil {
		db.Create(&models.App{
			Name:   "BizMart_service",
			Secret: "No secret in DB",
		})
	}

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) SaveUser(ctx context.Context, userRequest *ssov1.RegisterRequest) (models.User, error) {
	const op = "storage.sqlite.SaveUser"

	user := models.User{
		FirstName:    userRequest.GetFirstName(),
		LastName:     userRequest.GetLastName(),
		Email:        userRequest.GetEmail(),
		Username:     userRequest.GetUsername(),
		HashPassword: userRequest.GetPassword(),
	}
	result := s.db.WithContext(ctx).Create(&user)

	if result.Error != nil {
		// Проверяем, является ли ошибка sqlite3.Error и содержит ли она код UNIQUE_CONSTRAINT
		var sqliteErr sqlite3.Error
		if errors.As(result.Error, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			s.log.Warn("User already exists", slog.String("err", result.Error.Error()))
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		s.log.Error("Error inserting user", slog.String("err", result.Error.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, result.Error)
	}

	return user, nil
}

func (s *Storage) User(ctx context.Context, username string) (models.User, error) {
	const op = "storage.sqlite.User"

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
	const op = "storage.sqlite.IsAdmin"

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
	const op = "storage.sqlite.App"

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
