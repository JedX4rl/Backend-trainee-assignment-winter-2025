package storage

import (
	cfg "Backend-trainee-assignment-winter-2025/internal/config/storageConfig"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
	"time"
)

func NewStorage(dbConfig *cfg.StorageConfig) (*sql.DB, error) {

	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbConfig.Username, dbConfig.Password,
		dbConfig.Host, dbConfig.Port,
		dbConfig.DBName, dbConfig.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(12)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(20 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)
	applyMigrations(connectionString)

	return db, nil
}

func CallMigrations(dsn string) (*migrate.Migrate, error) {
	migrationsPath := "file://migrations"

	m, err := migrate.New(migrationsPath, dsn)

	if err != nil {
		return nil, err
	}

	return m, err
}

func applyMigrations(connStr string) {
	m, err := CallMigrations(connStr)
	if err != nil {
		log.Fatalf("Error initializing migrations: %v", err)
	}

	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			log.Println("Migrations are already applied")
		} else {
			log.Fatalf("Error applying migrations: %v", err)
		}
	} else {
		log.Println("Migrations successfully applied")
	}
}
