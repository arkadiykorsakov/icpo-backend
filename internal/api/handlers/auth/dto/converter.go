package dto

var allowedFields = map[string]string{
	"Email":    "email",
	"Password": "password",
	"Name":     "name",
}

func ToValidationErrorsResponse(errs map[string][]string) map[string][]string {
	resp := make(map[string][]string)
	for field, msgs := range errs {
		if key, ok := allowedFields[field]; ok {
			resp[key] = msgs
		}
	}
	return resp
}
