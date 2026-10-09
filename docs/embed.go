// Package docs содержит OpenAPI-спецификацию API, встроенную в бинарник.
package docs

import _ "embed"

// OpenAPI — содержимое openapi.yaml.
//
//go:embed openapi.yaml
var OpenAPI []byte
