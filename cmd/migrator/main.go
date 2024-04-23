package main

import (
	"auth_service/internal/config"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// go run ./cmd/migrator  --migrations-table=migrations --config-path=config/config_local.yaml --migrations-path=migrations
// go run ./cmd/migrator  --migrations-table=migrations_test --config-path=config/config_local.yaml --migrations-path=tests/migrations
func main() {
	var configPath, migrationsTable, migrationsPath string

	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "path to migrations table")

	flag.StringVar(&configPath, "config-path", "config", "path to config")

	flag.StringVar(&migrationsPath, "migrations-path", "migrations", "path to migrations")

	flag.Parse()

	cfg := config.MustLoadFromPath(configPath)

	if migrationsPath == "" {
		migrationsPath = cfg.MigrationPath
		if migrationsPath == "" {
			panic("migrations path is empty in config")
		}
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s&x-migrations-table=%s",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DBName,
		cfg.Postgres.SSLMode,
		migrationsTable,
	)

	m, err := migrate.New("file://"+migrationsPath, connStr)
	if err != nil {
		panic(fmt.Errorf("failed to initialize migrate: %w", err))
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
			return
		}
		panic(fmt.Errorf("migration failed: %w", err))
	}

	fmt.Println("Migrations applied successfully.")
}
