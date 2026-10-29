package city

import (
	"context"
	"fmt"
	"strings"

	"icpo-backend/internal/domain/city/dto"
	errs "icpo-backend/internal/errors"
	rootdto "icpo-backend/internal/shared/dto"
	"icpo-backend/internal/utils"
	"icpo-backend/pkg/pgerr"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	autocompleteLimit = 10

	colCityID        = "c.id"
	colCityName      = "c.name"
	colCityCreatedAt = "c.created_at"
	colCityCountryID = "c.country_id"

	defaultSortField = colCityName
	defaultSortDir   = "ASC"

	// прописано в миграции
	uniqueNameCityIndex = "cities_name_ci"
)

var allowedSortFields = map[string]string{
	"id":         colCityID,
	"name":       colCityName,
	"created_at": colCityCreatedAt,
	"country_id": colCityCountryID,
}

type querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *repository {
	return &repository{
		db: db,
	}
}

func (r *repository) list(
	ctx context.Context,
	q querier,
	search string,
	filter dto.FilterParams,
	orderBy string,
	prefixFirst bool,
	limit, offset int64,
) ([]City, error) {
	var (
		args       []any
		conditions []string
	)

	if search != "" {
		esc := utils.EscapeLike(search)

		args = append(args, "%"+esc+"%")
		conditions = append(
			conditions,
			fmt.Sprintf("c.name ILIKE $%d", len(args)),
		)

		if prefixFirst {
			args = append(args, esc+"%")
			orderBy = fmt.Sprintf(
				"CASE WHEN c.name ILIKE $%d THEN 0 ELSE 1 END, %s",
				len(args),
				orderBy,
			)
		}
	}

	if filter.CountryID != 0 {
		args = append(args, filter.CountryID)
		conditions = append(
			conditions,
			fmt.Sprintf("%s = $%d", colCityCountryID, len(args)),
		)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	args = append(args, limit, offset)
	limitParam := len(args) - 1
	offsetParam := len(args)

	query := fmt.Sprintf(`
        SELECT
            c.id,
            c.name,
            c.country_id,
            c.created_at,
            co.name       AS country_name,
            co.created_at AS country_created_at
        FROM cities c
        JOIN countries co ON co.id = c.country_id
        %s
        ORDER BY %s
        LIMIT $%d OFFSET $%d`,
		whereClause,
		orderBy,
		limitParam,
		offsetParam,
	)

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cities := make([]City, 0)
	for rows.Next() {
		var c City

		if err := rows.Scan(
			&c.ID,
			&c.Name,
			&c.Country.ID,
			&c.CreatedAt,
			&c.Country.Name,
			&c.Country.CreatedAt,
		); err != nil {
			return nil, err
		}

		cities = append(cities, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return cities, nil
}

func (r *repository) count(
	ctx context.Context,
	q querier,
	search string,
	filter dto.FilterParams,
) (int, error) {
	var (
		args       []any
		conditions []string
	)

	if search != "" {
		esc := utils.EscapeLike(search)

		args = append(args, "%"+esc+"%")
		conditions = append(
			conditions,
			fmt.Sprintf("%s ILIKE $%d", colCityName, len(args)),
		)
	}

	if filter.CountryID != 0 {
		args = append(args, filter.CountryID)
		conditions = append(
			conditions,
			fmt.Sprintf("%s = $%d", colCityCountryID, len(args)),
		)
	}

	query := "SELECT COUNT(*) FROM cities c"
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	if err := q.QueryRow(ctx, query, args...).Scan(&total); err != nil {
		return 0, err
	}

	return total, nil
}

func (r *repository) GetList(
	ctx context.Context,
	p rootdto.ListParams[dto.FilterParams],
) ([]City, int, error) {
	limit := int64(p.Pagination.Limit())
	offset := int64(p.Pagination.Offset())

	sortField := defaultSortField
	if p.Sort.Field != "" {
		var ok bool
		sortField, ok = allowedSortFields[p.Sort.Field]
		if !ok {
			return nil, 0, fmt.Errorf(
				"%w: %q",
				errs.ErrWrongSortField,
				p.Sort.Field,
			)
		}
	}

	sortDir := defaultSortDir
	if strings.EqualFold(string(p.Sort.Dir), "desc") {
		sortDir = "DESC"
	}

	orderBy := fmt.Sprintf("%s %s", sortField, sortDir)
	if sortField != colCityID {
		orderBy += ", " + colCityID + " ASC"
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel:   pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
	})
	if err != nil {
		return nil, 0, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	total, err := r.count(ctx, tx, p.Search, p.Filter)
	if err != nil {
		return nil, 0, err
	}

	items, err := r.list(
		ctx,
		tx,
		p.Search,
		p.Filter,
		orderBy,
		false,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, 0, err
	}

	return items, total, nil
}

func (r *repository) ExistsByName(ctx context.Context, name string, countryID int64) (bool, error) {
	const query = `
		SELECT EXISTS (SELECT 1 FROM cities WHERE LOWER(name) = LOWER($1) AND country_id = $2)
		`

	var exists bool
	if err := r.db.QueryRow(ctx, query, name, countryID).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *repository) Autocomplete(ctx context.Context, search string, filter dto.FilterParams) ([]City, error) {
	orderBy := colCityName + " ASC, " + colCityID + " ASC"
	return r.list(ctx, r.db, search, filter, orderBy, true, autocompleteLimit, 0)
}

func (r *repository) Create(ctx context.Context, name string, countryID int64) (City, error) {
	const query = `
		WITH inserted AS (
			INSERT INTO cities (name, country_id) VALUES ($1, $2)
			RETURNING id, name, country_id, created_at
		)
		SELECT i.id, i.name, i.country_id, i.created_at, co.name, co.created_at
		FROM inserted i JOIN countries co ON co.id = i.country_id`

	c := City{}
	err := r.db.QueryRow(ctx, query, name, countryID).Scan(
		&c.ID,
		&c.Name,
		&c.Country.ID,
		&c.CreatedAt,
		&c.Country.Name,
		&c.Country.CreatedAt,
	)

	if err != nil {
		if pgerr.IsUniqueViolation(err, uniqueNameCityIndex) {
			return City{}, errs.ErrCityExists
		}
		if pgerr.IsForeignKeyViolation(err) {
			return City{}, errs.ErrCountryNotFound
		}
		return City{}, err
	}

	return c, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM cities WHERE id = $1`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if pgerr.IsForeignKeyViolation(err) {
			return errs.ErrCityInUse
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return errs.ErrCityNotFound
	}
	return nil
}
