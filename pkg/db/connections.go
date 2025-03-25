package db

import (
	"SSO/internal/config"
	"SSO/pkg/db/postgres"
	sqliteDB "SSO/pkg/db/sqlite"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func ConnectToDB(cfg *config.Config) error {
	switch cfg.AppParams.DBSM {
	case "sqlite":
		err := sqliteDB.ConnectToDB(cfg.SqliteParams)
		if err != nil {
			return err
		}
	case "postgres":
		err := postgres.ConnectToDB(cfg.PostgresParams)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unsupported database: %s", cfg.AppParams.DBSM))
	}

	return nil
}

func GetDBConn(dbsm string) *gorm.DB {
	switch dbsm {
	case "sqlite":
		return sqliteDB.GetDBConn()
	case "postgres":
		return postgres.GetDBConn()
	default:
		fmt.Println(errors.New(fmt.Sprintf("unsupported database: %s", dbsm)))
	}

	return nil
}

func CloseDB(dbsm string) error {
	switch dbsm {
	case "sqlite":
		err := sqliteDB.CloseDBConn()
		if err != nil {
			return err
		}
	case "postgres":
		err := postgres.CloseDBConn()
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("unsupported database: %s", dbsm))
	}

	return nil
}
