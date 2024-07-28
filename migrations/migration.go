package migration

import (
	"log"
	"path/filepath"

	_ "github.com/lib/pq"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func MigrateUp() {
	absPath, err := filepath.Abs("./migrations/")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	m, err := migrate.New(
		"file://"+absPath,
		"postgres://postgres:root@localhost:5431/personal_ai?sslmode=disable")

	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}

func MigrateDown() {
	absPath, err := filepath.Abs("./migrations/")
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}

	m, err := migrate.New(
		"file://"+absPath,
		"postgres://postgres:root@localhost:5431/personal_ai?sslmode=disable")

	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
