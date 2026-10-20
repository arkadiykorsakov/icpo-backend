package errors

import "errors"

var (
	ErrEmailTaken         = errors.New("Email already taken")
	ErrInvalidCredentials = errors.New("Invalid email or password")
	ErrWrongPassword      = errors.New("Wrong password")

	ErrUserNotFound = errors.New("User not found")

	ErrInvalidToken = errors.New("Invalid token")

	ErrUserIDMissing = errors.New("userID missing in context after Authenticate")
)
