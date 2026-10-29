package dto

import (
	"icpo-backend/pkg/pagination"
	"icpo-backend/pkg/sorting"
)

type ListParams[T any] struct {
	Search     string
	Filter     T
	Pagination pagination.Params
	Sort       sorting.Params
}
