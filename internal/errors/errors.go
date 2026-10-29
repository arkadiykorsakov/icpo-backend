package errors

import "errors"

var (
	ErrEmailTaken         = errors.New("Email already taken")
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrWrongPassword      = errors.New("Wrong password")

	ErrUserNotFound = errors.New("User not found")

	ErrRefDataNotFound = errors.New("Ref data not found")
	ErrRefDataExists   = errors.New("Ref data already exists")
	ErrRefDataInUse    = errors.New("Ref data is in use")

	ErrCityNotFound = errors.New("City not found")
	ErrCityExists   = errors.New("City already exists")
	ErrCityInUse    = errors.New("City is in use")

	ErrCountryNotFound = errors.New("Country not found")

	ErrInvalidToken = errors.New("Invalid token")

	ErrUserIDMissing = errors.New("userID missing in context after Authenticate")

	ErrWrongSortField = errors.New("Wrong sort field")
)
