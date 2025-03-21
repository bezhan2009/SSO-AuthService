package main

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3" // Импорт драйвера SQLite
	_ "github.com/golang-migrate/migrate/v4/source/file"      // Импорт драйвера file
)

func main() {
	var storagePath, migrationsPath, migrationTable string

	flag.StringVar(&storagePath, "storage-path", "", "Path to a file containing the migration files")
	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a file containing the migration files")
	flag.StringVar(&migrationTable, "migration-table", "migrations", "Name of the migration table")
	flag.Parse()

	if storagePath == "" {
		log.Fatal("Missing storage path")
	}

	if migrationsPath == "" {
		log.Fatal("Missing migrations path")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationTable),
	)
	if err != nil {
		log.Fatal(err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	fmt.Println("Migrations applied successfully")
}
