package eyecolor

import (
	eyeColorDomain "icpo-backend/internal/domain/eye-color"
	"icpo-backend/internal/shared/refdata"
)

type Service interface {
	refdata.Service[eyeColorDomain.EyeColor]
}

// handler — HTTP-методы (GetAll, Autocomplete, Create, Delete) берутся из refdata.
type handler struct {
	refdata.Handler[eyeColorDomain.EyeColor]
}

func NewHandler(service Service) *handler {
	return &handler{
		Handler: refdata.NewHandler[eyeColorDomain.EyeColor](service),
	}
}
