package refdata

import (
	"context"
	"icpo-backend/internal/shared/refdata/dto"
)

type Repository[T any] interface {
	GetList(ctx context.Context, p dto.ListParams) ([]T, int, error)
	ExistsByName(ctx context.Context, name string) (bool, error)
	Autocomplete(ctx context.Context, search string) ([]T, error)
	Create(ctx context.Context, name string) (T, error)
	Delete(ctx context.Context, id int64) error
}

type service[T any] struct {
	repository Repository[T]
}

func NewService[T any](repository Repository[T]) *service[T] {
	return &service[T]{repository: repository}
}

func (s *service[T]) GetList(ctx context.Context, p dto.ListParams) ([]T, int, error) {
	return s.repository.GetList(ctx, p)
}

func (s *service[T]) Autocomplete(ctx context.Context, search string) ([]T, error) {
	return s.repository.Autocomplete(ctx, search)
}

func (s *service[T]) Create(ctx context.Context, name string) (T, error) {
	return s.repository.Create(ctx, name)
}

func (s *service[T]) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}

func (s *service[T]) ExistsByName(ctx context.Context, name string) (bool, error) {
	return s.repository.ExistsByName(ctx, name)
}
