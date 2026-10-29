package country

import (
	"context"

	"icpo-backend/internal/shared/refdata"
)

type Repository interface {
	refdata.Repository[Country]
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

type service struct {
	refdata.Service[Country]
}

func NewService(repository Repository) *service {
	return &service{
		Service: refdata.NewService[Country](repository),
	}
}
