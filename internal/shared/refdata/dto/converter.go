package dto

import "reflect"

func ToAutocompleteResponse[T any](p T) AutocompleteResponse {
	v := reflect.ValueOf(p)
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return AutocompleteResponse{}
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return AutocompleteResponse{}
	}

	var resp AutocompleteResponse
	if f := v.FieldByName("ID"); f.IsValid() && f.CanInt() {
		resp.ID = f.Int()
	}
	if f := v.FieldByName("Name"); f.IsValid() && f.Kind() == reflect.String {
		resp.Label = f.String()
	}
	return resp
}

func ToAutocompleteResponses[T any](items []T) []AutocompleteResponse {
	responses := make([]AutocompleteResponse, 0, len(items))
	for _, p := range items {
		responses = append(responses, ToAutocompleteResponse[T](p))
	}
	return responses
}
