package sorting

import "fmt"

type Direction string

const (
	Asc  Direction = "asc"
	Desc Direction = "desc"
)

type Params struct {
	Dir   Direction
	Field string
}

func Parse(fieldStr, dirStr string) (Params, error) {
	// Сортировка не задана вообще — это валидно.
	if fieldStr == "" && dirStr == "" {
		return Params{}, nil
	}

	// Направление без поля — противоречие.
	if fieldStr == "" {
		return Params{}, fmt.Errorf("sort field must be set when sort direction is provided")
	}

	dir := Direction(dirStr)
	if dir == "" {
		dir = Asc // дефолт, если поле задано, а направление нет
	}

	switch dir {
	case Asc, Desc:
	default:
		return Params{}, fmt.Errorf("invalid sort direction %q: must be %q or %q", dir, Asc, Desc)
	}

	return Params{Field: fieldStr, Dir: dir}, nil
}
