package utils

import "strings"

func NormalizeEmail(e string) string {
	return strings.ToLower(strings.TrimSpace(e))
}
