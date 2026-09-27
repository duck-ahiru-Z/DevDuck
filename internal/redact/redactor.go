package redact

import "regexp"

type Redactor struct {
	patterns []*regexp.Regexp
}

func New() *Redactor {
	return &Redactor{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)https?://[^\s/@:]+:[^\s/@]+@`),
			regexp.MustCompile(
				`(?i)(api[_-]?key|token|password|secret)\s*[:=]\s*["']?[^\s"']+`,
			),
			regexp.MustCompile(
				`(?i)authorization:\s*bearer\s+[^\s]+`,
			),
			regexp.MustCompile(
				`(?i)bearer\s+[A-Za-z0-9._~+/=-]+`,
			),
		},
	}
}

// Sanitize はAIへ送る前に秘密情報らしき文字列を隠す。
func (r *Redactor) Sanitize(text string) string {
	result := text

	for _, pattern := range r.patterns {
		result = pattern.ReplaceAllString(
			result,
			"[REDACTED]",
		)
	}

	return result
}
