package city

import (
	"context"
	"strings"

	"icpo-backend/internal/domain/city/dto"
	errs "icpo-backend/internal/errors"
	rootdto "icpo-backend/internal/shared/dto"
)

type CountryRepository interface {
	ExistsByID(ctx context.Context, id int64) (bool, error)
}

type CityRepository interface {
	GetList(ctx context.Context, p rootdto.ListParams[dto.FilterParams]) ([]City, int, error)
	ExistsByName(ctx context.Context, name string, countryID int64) (bool, error)
	Autocomplete(ctx context.Context, search string, filter dto.FilterParams) ([]City, error)
	Create(ctx context.Context, name string, countryID int64) (City, error)
	Delete(ctx context.Context, id int64) error
}

type service struct {
	countryRepo CountryRepository
	cityRepo    CityRepository
}

func NewService(cityRepo CityRepository, countryRepo CountryRepository) *service {
	return &service{
		countryRepo: countryRepo,
		cityRepo:    cityRepo,
	}
}

func (s *service) GetList(ctx context.Context, p rootdto.ListParams[dto.FilterParams]) ([]City, int, error) {
	return s.cityRepo.GetList(ctx, p)
}

func (s *service) Autocomplete(ctx context.Context, search string, filter dto.FilterParams) ([]City, error) {
	return s.cityRepo.Autocomplete(ctx, search, filter)
}

func (s *service) Create(ctx context.Context, name string, countryID int64) (City, error) {
	name = strings.TrimSpace(name)

	existsCountryID, err := s.countryRepo.ExistsByID(ctx, countryID)
	if err != nil {
		return City{}, err
	}
	if !existsCountryID {
		return City{}, errs.ErrCountryNotFound
	}

	existsCityName, err := s.cityRepo.ExistsByName(ctx, name, countryID)
	if err != nil {
		return City{}, err
	}
	if existsCityName {
		return City{}, errs.ErrCityExists
	}

	return s.cityRepo.Create(ctx, name, countryID)
}

func (s *service) Delete(ctx context.Context, id int64) error {
	return s.cityRepo.Delete(ctx, id)
}
