package refdata

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"icpo-backend/internal/api/httpx"
	errs "icpo-backend/internal/errors"
	"icpo-backend/internal/shared/refdata/dto"
	"icpo-backend/pkg/pagination"
	"icpo-backend/pkg/sorting"
	"icpo-backend/pkg/validator"
)

type Service[T any] interface {
	GetList(ctx context.Context, p dto.ListParams) ([]T, int, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Autocomplete(ctx context.Context, search string) ([]T, error)
	Create(ctx context.Context, name string) (T, error)
	Delete(ctx context.Context, id int64) error
}

type Handler[T any] interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Autocomplete(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type handler[T any] struct {
	service Service[T]
}

func NewHandler[T any](service Service[T]) *handler[T] {
	return &handler[T]{service: service}
}

// GetAll — GET ?page=1&per_page=20[&search=..][&sort=name|id|created_at][&dir=asc|desc]
func (h *handler[T]) GetAll(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	paginationParams, err := pagination.Parse(q.Get("page"), q.Get("per_page"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	sortParams, err := sorting.Parse(q.Get("sort"), q.Get("dir"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	if sortParams.Field != "" && !IsSortableField(sortParams.Field) {
		httpx.Error(w, http.StatusBadRequest, "unknown sort field "+sortParams.Field)
		return
	}

	listParams := dto.ListParams{
		Filter:     dto.FilterParams{Search: strings.TrimSpace(q.Get("search"))},
		Pagination: paginationParams,
		Sort:       sortParams,
	}

	list, total, err := h.service.GetList(r.Context(), listParams)
	if err != nil {
		httpx.ServerError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, pagination.Result[T]{
		Items:   list,
		Total:   total,
		Page:    paginationParams.Page,
		PerPage: paginationParams.PerPage,
	})
}

// Autocomplete — GET ?search=..
func (h *handler[T]) Autocomplete(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	items, err := h.service.Autocomplete(r.Context(), search)
	if err != nil {
		httpx.ServerError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"items": dto.ToAutocompleteResponses(items),
	})
}

func (h *handler[T]) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateRequest
	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if verrs := validator.ValidateRequest(input); verrs != nil {
		httpx.ValidationErrors(w, verrs)
		return
	}

	name := strings.ToLower(strings.TrimSpace(input.Name))

	exists, err := h.service.ExistsByName(r.Context(), name)
	if err != nil && !errors.Is(err, errs.ErrRefDataNotFound) {
		httpx.ServerError(w, err)
		return
	}

	if exists {
		httpx.ValidationErrors(w, map[string][]string{
			"name": {errs.ErrRefDataExists.Error()},
		})
	}

	item, err := h.service.Create(r.Context(), name)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrRefDataExists):
			httpx.ValidationErrors(w, map[string][]string{
				"name": {err.Error()},
			})
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, item)
}

func (h *handler[T]) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.GetIDFromRequest(w, r)
	if err != nil {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, errs.ErrRefDataNotFound):
			httpx.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, errs.ErrRefDataInUse):
			httpx.Error(w, http.StatusConflict, err.Error())
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
