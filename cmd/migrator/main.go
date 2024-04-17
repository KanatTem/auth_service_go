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

func main() {
	var cfgPath, migrationsTable string

	flag.StringVar(&cfgPath, "config", "config_local.yaml", "path to YAML config")

	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "path to migrations table")

	flag.Parse()

	cfg := config.MustLoad()

	migrationsPath := cfg.MigrationPath
	if migrationsPath == "" {
		panic("migrations path is empty in config")
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
