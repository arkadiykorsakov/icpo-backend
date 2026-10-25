package dto

import (
	"icpo-backend/pkg/pagination"
	"icpo-backend/pkg/sorting"
)

type FilterParams struct {
	Search string
}

type ListParams struct {
	Filter     FilterParams
	Pagination pagination.Params
	Sort       sorting.Params
}
