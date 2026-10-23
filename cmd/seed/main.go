// Команда seed загружает начальные данные справочников.
//
//	go run ./cmd/seed             # все справочники
//	go run ./cmd/seed countries   # только указанные
//
// Схема БД должна быть уже создана (миграции применяются при старте cmd/api).
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"icpo-backend/internal/config"
	"icpo-backend/internal/infrastructure/database"
	"icpo-backend/internal/infrastructure/database/seeder"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := database.NewPool(ctx, config.LoadDatabaseURL())
	if err != nil {
		return err
	}
	defer pool.Close()

	return seeder.Run(ctx, pool, os.Args[1:])
}
