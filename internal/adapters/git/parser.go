package git

import (
	"regexp"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

var firstDiagnostic = regexp.MustCompile(`(?m)^(?:fatal|error):\s+(.+)$`)

func Parse(stderr string) (model.ErrorInfo, bool) {
	lower := strings.ToLower(stderr)
	kind := "GitError"
	message := firstMessage(stderr)
	switch {
	case strings.Contains(stderr, "CONFLICT") || strings.Contains(lower, "merge conflict"):
		kind, message = "MergeConflict", conflictMessage(stderr)
	case strings.Contains(lower, "authentication failed") || strings.Contains(lower, "permission denied (publickey)"):
		kind = "AuthenticationError"
	case strings.Contains(lower, "could not resolve host") || strings.Contains(lower, "connection timed out") || strings.Contains(lower, "failed to connect"):
		kind = "NetworkError"
	case strings.Contains(lower, "not a git repository") || strings.Contains(lower, "dubious ownership"):
		kind = "RepositoryError"
	case strings.Contains(lower, "pathspec"):
		kind = "RefError"
	case strings.Contains(lower, "non-fast-forward") || strings.Contains(lower, "remote contains work that you do not have locally"):
		kind, message = "RefError", relevantMessage(stderr)
	case strings.Contains(lower, "local changes") && strings.Contains(lower, "overwritten"):
		kind = "WorktreeError"
	case strings.Contains(lower, "already exists") || strings.Contains(lower, "refspec"):
		kind = "RefError"
	case strings.Contains(lower, "unrelated histories"):
		kind = "RefError"
	case strings.Contains(lower, "please tell me who you are") || strings.Contains(lower, "user.name") || strings.Contains(lower, "user.email"):
		kind = "ConfigError"
	}
	if message == "" {
		return model.ErrorInfo{}, false
	}
	return model.ErrorInfo{Source: "git", Kind: kind, Message: message, Raw: stderr}, true
}

func firstMessage(stderr string) string {
	match := firstDiagnostic.FindStringSubmatch(stderr)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return relevantMessage(stderr)
}
func relevantMessage(stderr string) string {
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "hint:") {
			continue
		}
		if strings.Contains(strings.ToLower(line), "non-fast-forward") || strings.Contains(strings.ToLower(line), "remote contains work") || strings.Contains(line, "CONFLICT") || strings.HasPrefix(line, "!") {
			return line
		}
	}
	return firstMessageWithoutLoop(stderr)
}
func firstMessageWithoutLoop(stderr string) string {
	match := firstDiagnostic.FindStringSubmatch(stderr)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return strings.TrimSpace(strings.Split(stderr, "\n")[0])
}
func conflictMessage(stderr string) string { return relevantMessage(stderr) }
