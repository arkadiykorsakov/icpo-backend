package dto

import (
	"reflect"

	"icpo-backend/pkg/autocomplete"
)

func ToAutocompleteResponse[T any](p T) autocomplete.AutocompleteResult {
	v := reflect.ValueOf(p)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return autocomplete.AutocompleteResult{}
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return autocomplete.AutocompleteResult{}
	}

	var resp autocomplete.AutocompleteResult
	if f := v.FieldByName("ID"); f.IsValid() && f.CanInt() {
		resp.ID = f.Int()
	}
	if f := v.FieldByName("Name"); f.IsValid() && f.Kind() == reflect.String {
		resp.Label = f.String()
	}
	return resp
}
