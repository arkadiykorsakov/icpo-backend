package dto

import (
	"errors"
	"icpo-backend/internal/api/httpx"
	"net/http"

	cityDto "icpo-backend/internal/domain/city/dto"
)

var (
	ErrInvalidCountryID = errors.New("Invalid country id")
)

func ParseFilterParams(r *http.Request) (cityDto.FilterParams, error) {
	var f cityDto.FilterParams

	countryID, ok, err := httpx.ParseQueryNumber(r, "country_id", false)
	if err != nil {
		return f, err
	}
	if ok {
		if countryID == 0 {
			return f, ErrInvalidCountryID
		}
		f.CountryID = countryID
	}

	return f, nil
}
