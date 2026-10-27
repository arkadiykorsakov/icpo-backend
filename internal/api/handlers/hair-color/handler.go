package haircolor

import (
	hairColorDomain "icpo-backend/internal/domain/hair-color"
	"icpo-backend/internal/shared/refdata"
)

type Service interface {
	refdata.Service[hairColorDomain.HairColor]
}

// handler — HTTP-методы (GetAll, Autocomplete, Create, Delete) берутся из refdata.
type handler struct {
	refdata.Handler[hairColorDomain.HairColor]
}

func NewHandler(service Service) *handler {
	return &handler{
		Handler: refdata.NewHandler[hairColorDomain.HairColor](service),
	}
}
