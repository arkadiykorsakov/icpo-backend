package errors

import "errors"

var (
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrWrongPassword      = errors.New("wrong password")

	ErrUserNotFound = errors.New("user not found")

	ErrRefDataNotFound = errors.New("ref data not found")
	ErrRefDataExists   = errors.New("ref data already exists")
	ErrRefDataInUse    = errors.New("ref data is in use")

	ErrCityNotFound = errors.New("city not found")
	ErrCityExists   = errors.New("city already exists")
	ErrCityInUse    = errors.New("city is in use")

	ErrCountryNotFound = errors.New("country not found")

	ErrInvalidToken = errors.New("invalid token")

	ErrUserIDMissing = errors.New("userID missing in context after Authenticate")

	ErrWrongSortField = errors.New("wrong sort field")
)
