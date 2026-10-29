package dto

import (
	cityDomain "icpo-backend/internal/domain/city"
	"icpo-backend/pkg/autocomplete"
)

func ToAutocompleteResponse(c cityDomain.City) autocomplete.AutocompleteResult {
	return autocomplete.AutocompleteResult{
		ID:    c.ID,
		Label: c.Name + ", " + c.Country.Name,
	}
}
