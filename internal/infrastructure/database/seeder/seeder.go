// Package seeder загружает начальные данные справочников из seeds/*.txt.
//
// В отличие от миграций, сиды не версионируются: запуск идемпотентен
// (ON CONFLICT DO NOTHING), существующие записи не изменяются.
package seeder

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"icpo-backend/seeds"
)

const maxNameLen = 255 // совпадает с max=255 в dto.CreateRequest

// Tables — справочные таблицы в порядке загрузки. Для каждой должен существовать
// файл seeds/<table>.txt. Таблицы с внешними ключами добавляйте после тех,
// на которые они ссылаются.
var Tables = []string{
	"countries",
}

// Run загружает сиды в одной транзакции. Если only пуст — загружаются все
// таблицы из Tables, иначе только перечисленные (в порядке Tables).
func Run(ctx context.Context, pool *pgxpool.Pool, only []string) error {
	tables, err := selectTables(only)
	if err != nil {
		return err
	}

	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("seeder: begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck // после Commit это no-op

	for _, table := range tables {
		names, err := load(table)
		if err != nil {
			return err
		}

		query := fmt.Sprintf(
			`INSERT INTO %s (name) SELECT unnest($1::text[]) ON CONFLICT DO NOTHING`,
			pgx.Identifier{table}.Sanitize(),
		)
		tag, err := tx.Exec(ctx, query, names)
		if err != nil {
			return fmt.Errorf("seeder: %s: %w", table, err)
		}

		inserted := int(tag.RowsAffected())
		slog.Info("seeded", "table", table, "inserted", inserted, "skipped", len(names)-inserted)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("seeder: commit: %w", err)
	}
	return nil
}

func selectTables(only []string) ([]string, error) {
	if len(only) == 0 {
		return Tables, nil
	}
	for _, t := range only {
		if !slices.Contains(Tables, t) {
			return nil, fmt.Errorf("seeder: unknown table %q, available: %s", t, strings.Join(Tables, ", "))
		}
	}
	selected := make([]string, 0, len(only))
	for _, t := range Tables {
		if slices.Contains(only, t) {
			selected = append(selected, t)
		}
	}
	return selected, nil
}

// load читает seeds/<table>.txt: по названию на строку, «#» — комментарий.
func load(table string) ([]string, error) {
	file := table + ".txt"
	data, err := seeds.FS.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("seeder: read %s: %w", file, err)
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // BOM

	var names []string
	sc := bufio.NewScanner(bytes.NewReader(data))
	for line := 1; sc.Scan(); line++ {
		name := strings.TrimSpace(sc.Text())
		if name == "" || strings.HasPrefix(name, "#") {
			continue
		}
		if utf8.RuneCountInString(name) > maxNameLen {
			return nil, fmt.Errorf("seeder: %s:%d: name is longer than %d characters", file, line, maxNameLen)
		}
		names = append(names, name)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("seeder: scan %s: %w", file, err)
	}
	return names, nil
}
