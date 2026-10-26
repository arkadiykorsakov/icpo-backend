package language

import (
	"icpo-backend/internal/shared/refdata"
)

type Repository interface {
	refdata.Repository[Language]
}

type service struct {
	refdata.Service[Language]
}

func NewService(repository Repository) *service {
	return &service{
		Service: refdata.NewService[Language](repository),
	}
}
