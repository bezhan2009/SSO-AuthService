package postgres

import (
	"SSO/internal/config"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	dbConn *gorm.DB
)

func ConnectToDB(postgresParams config.PostgresParams) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		postgresParams.Host,
		postgresParams.Port,
		postgresParams.User,
		os.Getenv("DB_PASSWORD"),
		postgresParams.Database,
		postgresParams.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
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
