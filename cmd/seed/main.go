// Команда seed загружает начальные данные в базу.
//
//	go run ./cmd/seed             # все сидеры
//	go run ./cmd/seed NAME ...    # только указанные
//	go run ./cmd/seed -list       # показать доступные сидеры
//
// Схема БД должна быть уже создана (миграции применяются при старте cmd/api).
package main

import (
	"context"
	"flag"
	"fmt"
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
	list := flag.Bool("list", false, "показать доступные сидеры и выйти")
	flag.Parse()

	if *list {
		for _, name := range seeder.Names() {
			fmt.Println(name)
		}
		return nil
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	pool, err := database.NewPool(ctx, config.LoadDatabaseURL())
	if err != nil {
		return err
	}
	defer pool.Close()

	return seeder.Run(ctx, pool, flag.Args())
}
