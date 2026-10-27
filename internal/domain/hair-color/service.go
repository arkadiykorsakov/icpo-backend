package haircolor

import (
	"icpo-backend/internal/shared/refdata"
)

type Repository interface {
	refdata.Repository[HairColor]
}

type service struct {
	refdata.Service[HairColor]
}

func NewService(repository Repository) *service {
	return &service{
		Service: refdata.NewService[HairColor](repository),
	}
}
