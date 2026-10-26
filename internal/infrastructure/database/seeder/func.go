package seeder

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Func создаёт сидер из функции. Это основной способ написать сидер произвольной
// логики: функция получает транзакцию и сама решает, что и как загружать.
//
//	seeder.Func("admin", func(ctx context.Context, tx pgx.Tx) (seeder.Result, error) {
//		tag, err := tx.Exec(ctx, `INSERT ... ON CONFLICT DO NOTHING`)
//		if err != nil {
//			return seeder.Result{}, err
//		}
//		n := int(tag.RowsAffected())
//		return seeder.Result{Inserted: n, Skipped: 1 - n}, nil
//	})
func Func(name string, fn func(ctx context.Context, tx pgx.Tx) (Result, error)) Seeder {
	return funcSeeder{name: name, fn: fn}
}

type funcSeeder struct {
	name string
	fn   func(ctx context.Context, tx pgx.Tx) (Result, error)
}

func (s funcSeeder) Name() string { return s.name }

func (s funcSeeder) Run(ctx context.Context, tx pgx.Tx) (Result, error) {
	return s.fn(ctx, tx)
}
