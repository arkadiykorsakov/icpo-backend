// Package seeder загружает начальные данные в базу.
//
// В отличие от миграций, сидеры не версионируются и запускаются вручную
// командой cmd/seed. Каждый сидер обязан быть идемпотентным: повторный запуск
// не должен дублировать или портить данные.
//
// Все выбранные сидеры выполняются в одной транзакции в порядке, заданном
// в registry (registry.go): при ошибке любого из них не сохраняется ничего.
package seeder

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Seeder — единица начальных данных.
type Seeder interface {
	// Name — уникальное имя сидера: используется в командной строке и в логах.
	Name() string

	// Run выполняет сидер внутри транзакции tx и возвращает, сколько записей
	// было добавлено, а сколько пропущено как уже существующие.
	// Транзакцию (Commit/Rollback) Run не завершает — этим занимается раннер.
	Run(ctx context.Context, tx pgx.Tx) (Result, error)
}

// Result — итог работы одного сидера.
type Result struct {
	Inserted int // добавлено новых записей
	Skipped  int // пропущено (уже существовали)
}

// Run выполняет сидеры из registry в одной транзакции. Если only пуст —
// выполняются все, иначе только перечисленные по имени (в порядке registry).
func Run(ctx context.Context, pool *pgxpool.Pool, only []string) error {
	return run(ctx, pool, registry, only)
}

// Names возвращает имена всех зарегистрированных сидеров в порядке выполнения.
func Names() []string {
	return names(registry)
}

func run(ctx context.Context, pool *pgxpool.Pool, all []Seeder, only []string) error {
	selected, err := selectSeeders(all, only)
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("seeder: begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // после Commit это no-op

	for _, s := range selected {
		res, err := s.Run(ctx, tx)
		if err != nil {
			return fmt.Errorf("seeder: %s: %w", s.Name(), err)
		}
		slog.Info("seeded", "seeder", s.Name(), "inserted", res.Inserted, "skipped", res.Skipped)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("seeder: commit: %w", err)
	}
	return nil
}

// selectSeeders проверяет уникальность имён и отбирает сидеры по only,
// сохраняя порядок all.
func selectSeeders(all []Seeder, only []string) ([]Seeder, error) {
	index := make(map[string]struct{}, len(all))
	for _, s := range all {
		if _, dup := index[s.Name()]; dup {
			return nil, fmt.Errorf("seeder: duplicate name %q in registry", s.Name())
		}
		index[s.Name()] = struct{}{}
	}

	if len(only) == 0 {
		return all, nil
	}

	want := make(map[string]struct{}, len(only))
	for _, name := range only {
		if _, ok := index[name]; !ok {
			return nil, fmt.Errorf("seeder: unknown seeder %q, available: %s", name, strings.Join(names(all), ", "))
		}
		want[name] = struct{}{}
	}

	selected := make([]Seeder, 0, len(want))
	for _, s := range all {
		if _, ok := want[s.Name()]; ok {
			selected = append(selected, s)
		}
	}
	return selected, nil
}

func names(all []Seeder) []string {
	out := make([]string, 0, len(all))
	for _, s := range all {
		out = append(out, s.Name())
	}
	return out
}
