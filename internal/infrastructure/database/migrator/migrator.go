package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"icpo-backend/migrations"
)

func Up(ctx context.Context, databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("close db", "error", err)
		}
	}()

	provider, err := goose.NewProvider(goose.DialectPostgres, db, migrations.FS)
	if err != nil {
		return fmt.Errorf("goose provider: %w", err)
	}

	results, err := provider.Up(ctx)
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	if len(results) == 0 {
		slog.Info("migrations: no changes")
		return nil
	}

	for _, result := range results {
		slog.Info(
			"migration applied",
			"name", result.Source.Path,
			"duration", result.Duration,
		)
	}

	return nil
}
