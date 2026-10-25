package refdata

import (
	"context"
	"fmt"
	"strings"

	errs "icpo-backend/internal/errors"
	"icpo-backend/internal/shared/refdata/dto"
	"icpo-backend/internal/utils"
	"icpo-backend/pkg/pgerr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	autocompleteLimit = 10
	defaultSortField  = "name"
)

var sortableFields = map[string]struct{}{
	"id":         {},
	"name":       {},
	"created_at": {},
}

func IsSortableField(field string) bool {
	_, ok := sortableFields[field]
	return ok
}

type repository[T any] struct {
	db        *pgxpool.Pool
	tableName string
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewRepository[T any](db *pgxpool.Pool, tableName string) *repository[T] {
	return &repository[T]{db: db, tableName: tableName}
}

// table — безопасно экранированное имя таблицы.
func (r *repository[T]) table() string {
	return pgx.Identifier{r.tableName}.Sanitize()
}

func (r *repository[T]) list(ctx context.Context, q querier, search, orderBy string, prefixFirst bool, limit, offset int64) ([]T, error) {
	var (
		args   []any
		filter string
	)

	if search != "" {
		esc := utils.EscapeLike(search)
		args = append(args, "%"+esc+"%")
		filter = fmt.Sprintf("WHERE name ILIKE $%d", len(args))

		if prefixFirst {
			args = append(args, esc+"%")
			orderBy = fmt.Sprintf("CASE WHEN name ILIKE $%d THEN 0 ELSE 1 END, %s", len(args), orderBy)
		}
	}

	args = append(args, limit, offset)
	query := fmt.Sprintf(
		`SELECT id, name, created_at FROM %s %s ORDER BY %s LIMIT $%d OFFSET $%d`,
		r.table(), filter, orderBy, len(args)-1, len(args),
	)

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	ptrs, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[T])
	if err != nil {
		return nil, err
	}

	items := make([]T, len(ptrs))
	for i, p := range ptrs {
		items[i] = *p
	}
	return items, nil
}

func (r *repository[T]) count(ctx context.Context, q querier, search string) (int, error) {
	var (
		args   []any
		filter string
	)
	if search != "" {
		esc := utils.EscapeLike(search)
		args = append(args, "%"+esc+"%")
		filter = "WHERE name ILIKE $1"
	}

	query := fmt.Sprintf(`SELECT COUNT(*) FROM %s %s`, r.table(), filter)

	var total int

	if err := q.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *repository[T]) GetList(
	ctx context.Context,
	p dto.ListParams,
) ([]T, int, error) {
	limit := int64(p.Pagination.Limit())
	offset := int64(p.Pagination.Offset())

	field := defaultSortField
	if IsSortableField(p.Sort.Field) {
		field = p.Sort.Field
	}
	dir := "ASC"
	if strings.EqualFold(string(p.Sort.Dir), "desc") {
		dir = "DESC"
	}
	orderBy := fmt.Sprintf("%s %s", field, dir)
	if field != "id" {
		orderBy += ", id ASC"
	}

	search := p.Filter.Search

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx) // для read-only достаточно; после Commit это no-op

	total, err := r.count(ctx, tx, search)
	if err != nil {
		return nil, 0, err
	}

	items, err := r.list(ctx, tx, search, orderBy, false, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *repository[T]) Autocomplete(ctx context.Context, search string) ([]T, error) {
	return r.list(ctx, r.db, search, "name ASC", true, autocompleteLimit, 0)
}

func (r *repository[T]) Create(ctx context.Context, name string) (T, error) {
	var zero T

	query := fmt.Sprintf(
		`INSERT INTO %s (name) VALUES ($1) RETURNING id, name, created_at`,
		r.table(),
	)

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return zero, err
	}

	// Ошибки сервера (в т.ч. нарушение уникальности) приходят при чтении строк.
	item, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[T])
	if pgerr.IsUniqueViolation(err, "") {
		return zero, errs.ErrRefDataExists
	}
	if err != nil {
		return zero, err
	}
	return *item, nil
}

func (r *repository[T]) Delete(ctx context.Context, id int64) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1`, r.table())

	tag, err := r.db.Exec(ctx, query, id)
	if pgerr.IsForeignKeyViolation(err) {
		return errs.ErrRefDataInUse
	}
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrRefDataNotFound
	}
	return nil
}

func (r *repository[T]) ExistsByName(ctx context.Context, name string) (bool, error) {
	query := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM %s WHERE LOWER(name) = LOWER($1))`, r.table())

	var exists bool
	if err := r.db.QueryRow(ctx, query, name).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
