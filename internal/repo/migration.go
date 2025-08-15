package repo

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Nadzyoki/BannerCounter/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			cfg.Host,
			cfg.Port,
			cfg.User,
			cfg.DB,
			cfg.Password,
		),
	)
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error("create driver", "err", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		slog.Error("create migrate instance", "err", err)
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migrate up", "err", err)
		return err
	}

	slog.Info("Migrations applied")
	return nil
}
