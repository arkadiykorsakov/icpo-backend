package eyecolor

import (
	"icpo-backend/internal/shared/refdata"
)

type Repository interface {
	refdata.Repository[EyeColor]
}

type service struct {
	refdata.Service[EyeColor]
}

func NewService(repository Repository) *service {
	return &service{
		Service: refdata.NewService[EyeColor](repository),
	}
}
