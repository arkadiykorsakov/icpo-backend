package city

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"icpo-backend/internal/api/handlers/city/dto"
	"icpo-backend/internal/api/httpx"
	cityDomain "icpo-backend/internal/domain/city"
	cityDto "icpo-backend/internal/domain/city/dto"
	errs "icpo-backend/internal/errors"
	rootdto "icpo-backend/internal/shared/dto"
	"icpo-backend/pkg/autocomplete"
	"icpo-backend/pkg/pagination"
	"icpo-backend/pkg/sorting"
	"icpo-backend/pkg/validator"
)

type Service interface {
	GetList(ctx context.Context, p rootdto.ListParams[cityDto.FilterParams]) ([]cityDomain.City, int, error)
	Autocomplete(ctx context.Context, search string, filter cityDto.FilterParams) ([]cityDomain.City, error)
	Create(ctx context.Context, name string, countryID int64) (cityDomain.City, error)
	Delete(ctx context.Context, id int64) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
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

	listParams := rootdto.ListParams[cityDto.FilterParams]{
		Search:     strings.TrimSpace(q.Get("search")),
		Pagination: paginationParams,
		Sort:       sortParams,
	}

	filter, err := dto.ParseFilterParams(r)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	listParams.Filter = filter

	list, total, err := h.service.GetList(r.Context(), listParams)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrWrongSortField):
			httpx.Error(w, http.StatusBadRequest, err.Error())
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.JSON(w, http.StatusOK, pagination.Result[cityDomain.City]{
		Items:   list,
		Total:   total,
		Page:    paginationParams.Page,
		PerPage: paginationParams.PerPage,
	})
}

func (h *handler) Autocomplete(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	filter, err := dto.ParseFilterParams(r)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	items, err := h.service.Autocomplete(r.Context(), search, filter)
	if err != nil {
		httpx.ServerError(w, err)
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"items": autocomplete.ToAutocompleteResponses(items, dto.ToAutocompleteResponse),
	})
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateRequest
	if err := httpx.DecodeJSON(w, r, &input); err != nil {
		return
	}

	if verrs := validator.ValidateRequest(input); verrs != nil {
		httpx.ValidationErrors(w, verrs)
		return
	}

	item, err := h.service.Create(r.Context(), input.Name, input.CountryID)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrCityExists):
			httpx.ValidationErrors(w, map[string][]string{
				"name": {err.Error()},
			})
		case errors.Is(err, errs.ErrCountryNotFound):
			httpx.ValidationErrors(w, map[string][]string{
				"country_id": {err.Error()},
			})
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	httpx.JSON(w, http.StatusCreated, item)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpx.GetIDFromRequest(w, r)
	if err != nil {
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, errs.ErrCityNotFound):
			slog.Info("delete: city already absent",
				"err", err,
				"id", id,
			)
			w.WriteHeader(http.StatusNoContent)
		case errors.Is(err, errs.ErrCityInUse):
			httpx.Error(w, http.StatusConflict, err.Error())
		default:
			httpx.ServerError(w, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
