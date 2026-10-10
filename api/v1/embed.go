// Package v1 содержит OpenAPI-спецификацию API, встроенную в бинарник.
package v1

import _ "embed"

// OpenAPI — содержимое openapi.yaml.
//
//go:embed openapi.yaml
var OpenAPI []byte
