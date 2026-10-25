package country

import (
	"icpo-backend/internal/shared/refdata"
)

type Repository interface {
	refdata.Repository[Country]
}

type service struct {
	refdata.Service[Country]
}

func NewService(repository Repository) *service {
	return &service{
		Service: refdata.NewService[Country](repository),
	}
}
