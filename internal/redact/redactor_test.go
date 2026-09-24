package redact

import (
	"strings"
	"testing"
)

func TestSanitizeAPIKey(t *testing.T) {
	redactor := New()

	input := "API_KEY=very-secret-value"

	result := redactor.Sanitize(input)

	if strings.Contains(
		result,
		"very-secret-value",
	) {
		t.Fatal(
			"expected API key to be redacted",
		)
	}
}

func TestSanitizeBearerToken(t *testing.T) {
	redactor := New()

	input := "Authorization: Bearer abc123xyz"

	result := redactor.Sanitize(input)

	if strings.Contains(
		result,
		"abc123xyz",
	) {
		t.Fatal(
			"expected bearer token to be redacted",
		)
	}
}
