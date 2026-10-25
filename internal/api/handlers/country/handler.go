package country

import (
	countryDomain "icpo-backend/internal/domain/country"
	"icpo-backend/internal/shared/refdata"
)

type Service interface {
	refdata.Service[countryDomain.Country]
}

// handler — HTTP-методы (GetAll, Autocomplete, Create, Delete) берутся из refdata.
type handler struct {
	refdata.Handler[countryDomain.Country]
}

func NewHandler(service Service) *handler {
	return &handler{
		Handler: refdata.NewHandler[countryDomain.Country](service),
	}
}
