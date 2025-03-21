package sqlite

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log/slog"
)

type Storage struct {
	db  *sql.DB
	log *slog.Logger
}

// New creates a new instance of the SQLite storage.
func New(storagePath string, log *slog.Logger) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", op, err)
	}

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?,?)")
	if err != nil {
		s.log.Error("Error preparing request", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s : %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		//var sqliteErr sqlite3.Error
		//if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		//	s.log.Warn("Failed to insert user", slog.String("err", err.Error()))
		//	return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		//}
		s.log.Error("Error inserting user", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		s.log.Error("Error lastInsertId", slog.String("err", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil // Приводим int64 к uint
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT * FROM users WHERE email = ?")
	if err != nil {
		s.log.Error("Error while preparing statement", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close() // Закрываем statement

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash, &user.IsAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.log.Warn("User not found", slog.String("email", email))
			return models.User{}, storage.ErrUserNotFound
		}
		s.log.Error("Error while scanning row", slog.String("err", err.Error()))
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		s.log.Error("Error prepareStatement", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID) // Приводим uint64 к int64
	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		s.log.Error("Error row scan", slog.String("err", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT * FROM apps WHERE id = ?")
	if err != nil {
		s.log.Error("Error prepareStatement", slog.String("err", err.Error()))
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, appID) // Приводим uint64 к int64

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, storage.ErrAppNotFound
		}
		s.log.Error("Error row scan", slog.String("err", err.Error()))
	}

	return app, nil
}
