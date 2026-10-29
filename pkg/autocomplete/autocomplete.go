package autocomplete

type AutocompleteResult struct {
	ID    int64  `json:"id"`
	Label string `json:"label"`
}

type ToAutocompleteFunc[T any] func(T) AutocompleteResult

func ToAutocompleteResponses[T any](
	items []T,
	toAutocomplete ToAutocompleteFunc[T],
) []AutocompleteResult {
	responses := make([]AutocompleteResult, 0, len(items))

	for _, item := range items {
		responses = append(responses, toAutocomplete(item))
	}

	return responses
}
