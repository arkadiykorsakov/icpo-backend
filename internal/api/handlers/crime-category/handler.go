package crimecategory

import (
	crimeCategoryDomain "icpo-backend/internal/domain/crime-category"
	"icpo-backend/internal/shared/refdata"
)

type Service interface {
	refdata.Service[crimeCategoryDomain.CrimeCategory]
}

// handler — HTTP-методы (GetAll, Autocomplete, Create, Delete) берутся из refdata.
type handler struct {
	refdata.Handler[crimeCategoryDomain.CrimeCategory]
}

func NewHandler(service Service) *handler {
	return &handler{
		Handler: refdata.NewHandler[crimeCategoryDomain.CrimeCategory](service),
	}
}
