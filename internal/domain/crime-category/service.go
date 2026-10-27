package crimecategory

import (
	"icpo-backend/internal/shared/refdata"
)

type Repository interface {
	refdata.Repository[CrimeCategory]
}

type service struct {
	refdata.Service[CrimeCategory]
}

func NewService(repository Repository) *service {
	return &service{
		Service: refdata.NewService[CrimeCategory](repository),
	}
}
