package pagination

import (
	"fmt"
	"strconv"
	"strings"
)

type Params struct {
	PerPage int
	Page    int
}

func (p Params) Limit() int  { return p.PerPage }
func (p Params) Offset() int { return (p.Page - 1) * p.PerPage }

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MinPerPage     = 1
	MaxPerPage     = 100
	MaxPage        = 1_000_000 // защита от переполнения в Offset()
)

func Parse(pageStr, perPageStr string) (Params, error) {
	page, err := parseIntOr(pageStr, DefaultPage)
	if err != nil {
		return Params{}, fmt.Errorf("page must be an integer")
	}
	perPage, err := parseIntOr(perPageStr, DefaultPerPage)
	if err != nil {
		return Params{}, fmt.Errorf("per_page must be an integer")
	}

	if page < 1 || page > MaxPage {
		return Params{}, fmt.Errorf("page must be in [1,%d]", MaxPage)
	}
	if perPage < MinPerPage || perPage > MaxPerPage {
		return Params{}, fmt.Errorf("per_page must be in [%d,%d]", MinPerPage, MaxPerPage)
	}

	return Params{Page: page, PerPage: perPage}, nil
}

// parseIntOr возвращает def, если строка пустая (после обрезки пробелов).
func parseIntOr(s string, def int) (int, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return def, nil
	}
	return strconv.Atoi(s)
}

type Result[T any] struct {
	Items   []T `json:"items"`
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
