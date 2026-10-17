package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()

	// В ошибках используем имена полей из json-тегов (email), а не Go-имена (Email).
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// только латиница, цифры и спецсимволы
	var passwordRe = regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?~\x60]+$`)
	mustRegister(v, "passwordchars", func(fl validator.FieldLevel) bool {
		return passwordRe.MatchString(fl.Field().String())
	})

	// maxbytes=N — ограничение длины строки в БАЙТАХ. Нужно для паролей:
	// bcrypt принимает не более 72 байт, а для кириллицы это всего 36 символов.
	mustRegister(v, "maxbytes", func(fl validator.FieldLevel) bool {
		limit, err := strconv.Atoi(fl.Param())
		if err != nil {
			return false
		}
		return len(fl.Field().String()) <= limit
	})

	// notblank — строка не пустая и не состоит из одних пробелов.
	mustRegister(v, "notblank", func(fl validator.FieldLevel) bool {
		return strings.TrimSpace(fl.Field().String()) != ""
	})

	return v
}

func mustRegister(v *validator.Validate, tag string, fn validator.Func) {
	if err := v.RegisterValidation(tag, fn); err != nil {
		panic(fmt.Sprintf("validator: register %q: %v", tag, err))
	}
}

func ValidateRequest(req any) map[string][]string {
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
		case "required", "notblank":
			msg = fmt.Sprintf("Field %s is required", field)
		case "email":
			msg = fmt.Sprintf("Field %s must be a valid email", field)
		case "min":
			msg = fmt.Sprintf("Field %s must have at least %s characters", field, e.Param())
		case "max":
			msg = fmt.Sprintf("Field %s must have no more than %s characters", field, e.Param())
		case "maxbytes":
			msg = fmt.Sprintf("Field %s must be no longer than %s bytes", field, e.Param())
		case "passwordchars":
			msg = fmt.Sprintf("Field %s may contain only Latin letters, digits and special symbols", field)
		default:
			msg = fmt.Sprintf("Field %s is invalid", field)
		}

		messages[field] = append(messages[field], msg)
	}
	return messages
}
