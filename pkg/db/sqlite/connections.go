package sqliteDB

import (
	"SSO/internal/config"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"log/slog"
)

var (
	dbConn *gorm.DB
)

func ConnectToDB(sqliteParams config.SqliteParams) error {
	const op = "sqlite.ConnectToDB"

	db, err := gorm.Open(sqlite.Open(sqliteParams.StoragePath), &gorm.Config{})
	if err != nil {
		fmt.Println(fmt.Sprintf("op: %s: Error opening database", op), slog.String("error", err.Error()))
		return err
	}

	dbConn = db

	return nil
}

func GetDBConn() *gorm.DB {
	return dbConn
}

func CloseDBConn() error {
	if sqlDB, err := GetDBConn().DB(); err == nil {
		if err = sqlDB.Close(); err != nil {
			log.Fatalf("Error while closing DB: %s", err)
		}
		fmt.Println("Connection closed successfully")
	} else {
		log.Fatalf("Error while getting *sql.DB from GORM: %s", err)
	}

	return nil
}
