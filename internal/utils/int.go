package utils

import (
	"errors"
	"strconv"
)

var (
	ErrNotInt      = errors.New("must be an integer")
	ErrNotPositive = errors.New("must be a positive integer")
)

func ParsePositiveInt64(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, ErrNotInt
	}
	if n < 1 {
		return 0, ErrNotPositive
	}
	return n, nil
}
