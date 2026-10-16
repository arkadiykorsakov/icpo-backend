package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest(req interface{}) map[string][]string {
	err := validate.Struct(req)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return map[string][]string{
			"_error": {"Validation error not found"},
		}
	}

	messages := make(map[string][]string)
	for _, e := range validationErrors {
		field := e.Field()

		var msg string
		switch e.Tag() {
		case "required":
			msg = fmt.Sprintf("Field %s is required", field)
		case "email":
			msg = fmt.Sprintf("Field %s must be a valid email", field)
		case "min":
			msg = fmt.Sprintf("Field %s must have at least %s characters", field, e.Param())
		case "max":
			msg = fmt.Sprintf("Field %s must have no more than %s characters", field, e.Param())
		default:
			msg = fmt.Sprintf("Field %s is invalid", field)
		}

		messages[field] = append(messages[field], msg)
	}
	return messages
}
