package language

import (
	languageDomain "icpo-backend/internal/domain/language"
	"icpo-backend/internal/shared/refdata"
)

type Service interface {
	refdata.Service[languageDomain.Language]
}

// handler — HTTP-методы (GetAll, Autocomplete, Create, Delete) берутся из refdata.
type handler struct {
	refdata.Handler[languageDomain.Language]
}

func NewHandler(service Service) *handler {
	return &handler{
		Handler: refdata.NewHandler[languageDomain.Language](service),
	}
}
