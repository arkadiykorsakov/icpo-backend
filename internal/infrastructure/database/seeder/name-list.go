package seeder

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"

	"icpo-backend/seeds"
)

const maxNameLen = 255 // совпадает с max=255 в dto.CreateRequest

// NameList создаёт сидер, который читает seeds/<table>.txt (по значению на строку)
// и вставляет значения в колонку name таблицы <table>. Имя сидера совпадает
// с именем таблицы.
//
// Формат файла: пустые строки и строки, начинающиеся с «#», игнорируются,
// пробелы по краям обрезаются. Дубликаты отсекаются через ON CONFLICT DO NOTHING,
// поэтому в таблице должен быть уникальный индекс по name.
func NameList(table string) Seeder {
	return nameList{table: table}
}

type nameList struct{ table string }

func (s nameList) Name() string { return s.table }

func (s nameList) Run(ctx context.Context, tx pgx.Tx) (Result, error) {
	file := s.table + ".txt"
	data, err := seeds.FS.ReadFile(file)
	if err != nil {
		return Result{}, fmt.Errorf("read %s: %w", file, err)
	}
	values, err := parseNameList(file, data)
	if err != nil {
		return Result{}, err
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (name) SELECT unnest($1::text[]) ON CONFLICT DO NOTHING`,
		pgx.Identifier{s.table}.Sanitize(),
	)
	tag, err := tx.Exec(ctx, query, values)
	if err != nil {
		return Result{}, err
	}

	inserted := int(tag.RowsAffected())
	return Result{Inserted: inserted, Skipped: len(values) - inserted}, nil
}

// parseNameList разбирает содержимое файла со списком значений.
func parseNameList(file string, data []byte) ([]string, error) {
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")) // BOM

	var values []string
	sc := bufio.NewScanner(bytes.NewReader(data))
	for line := 1; sc.Scan(); line++ {
		v := strings.TrimSpace(sc.Text())
		if v == "" || strings.HasPrefix(v, "#") {
			continue
		}
		if utf8.RuneCountInString(v) > maxNameLen {
			return nil, fmt.Errorf("%s:%d: value is longer than %d characters", file, line, maxNameLen)
		}
		values = append(values, v)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", file, err)
	}
	return values, nil
}
