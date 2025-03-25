package app

import (
	grpcapp "SSO/internal/app/grpc"
	"SSO/internal/config"
	"SSO/internal/services/auth"
	"SSO/internal/storage"
	"SSO/internal/storage/postgres"
	"SSO/internal/storage/sqlite"
	dbPostgres "SSO/pkg/db/postgres"
	dbSqlite "SSO/pkg/db/sqlite"
	"fmt"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	var s storage.Storage
	var err error

	switch cfg.AppParams.DBSM {
	case "sqlite":
		s, err = sqlite.New(dbSqlite.GetDBConn(), log)
	case "postgres":
		s, err = postgres.New(dbPostgres.GetDBConn(), log)
	default:
		panic(fmt.Sprintf("unsupported database: %s", cfg.AppParams.DBSM))
	}

	if err != nil {
		panic(err)
	}

	authService := auth.New(log, s, s, s, cfg.AuthParams)

	grpcApp := grpcapp.New(log, authService, cfg.GRPC.Port)

	return &App{
		GRPCServer: grpcApp,
	}
}
