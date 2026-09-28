package python

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

var filePattern = regexp.MustCompile(
	`File "([^"]+)", line ([0-9]+)`,
)

var errorPattern = regexp.MustCompile(
	`(?m)^((?:[A-Za-z_][A-Za-z0-9_]*\.)*[A-Za-z_][A-Za-z0-9_]*)(?::\s*(.*))$`,
)

var controlPattern = regexp.MustCompile(`(?m)^(KeyboardInterrupt|SystemExit|GeneratorExit)$`)
var noMessagePattern = regexp.MustCompile(`^(?:[A-Za-z_][A-Za-z0-9_]*\.)*[A-Za-z_][A-Za-z0-9_]*$`)

func Parse(stderr string) (model.ErrorInfo, bool) {
	normalized := strings.ReplaceAll(stderr, "\r\n", "\n")
	errorMatches := errorPattern.FindAllStringSubmatch(normalized, -1)
	if len(errorMatches) == 0 {
		controlMatches := controlPattern.FindAllStringSubmatch(normalized, -1)
		if len(controlMatches) > 0 {
			return model.ErrorInfo{Source: "python", Kind: controlMatches[len(controlMatches)-1][1], Raw: stderr, SkipTeaching: true}, true
		}
		if strings.Contains(normalized, "Traceback") {
			lines := strings.Split(strings.TrimRight(normalized, "\r\n"), "\n")
			for index := len(lines) - 1; index >= 0; index-- {
				candidate := strings.TrimSpace(lines[index])
				if candidate == "" {
					continue
				}
				if noMessagePattern.MatchString(candidate) {
					return withLocation(parseInfo(candidate, "", stderr), normalized), true
				}
				break
			}
		}
	}

	if len(errorMatches) == 0 {
		return model.ErrorInfo{}, false
	}

	lastError := errorMatches[len(errorMatches)-1]

	info := parseInfo(lastError[1], lastError[2], stderr)

	info = withLocation(info, normalized)
	return info, true
}

func withLocation(info model.ErrorInfo, stderr string) model.ErrorInfo {
	fileMatches := filePattern.FindAllStringSubmatch(stderr, -1)

	if len(fileMatches) > 0 {
		lastFile := fileMatches[len(fileMatches)-1]

		info.File = lastFile[1]

		line, err := strconv.Atoi(lastFile[2])

		if err == nil {
			info.Line = line
		}
	}

	return info
}

func parseInfo(name, message, raw string) model.ErrorInfo {
	return model.ErrorInfo{Source: "python", Kind: unqualifiedKind(name), Message: message, Raw: raw, SkipTeaching: IsControlFlow(unqualifiedKind(name))}
}

func unqualifiedKind(name string) string {
	for index := len(name) - 1; index >= 0; index-- {
		if name[index] == '.' {
			return name[index+1:]
		}
	}
	return name
}

func IsControlFlow(kind string) bool {
	switch kind {
	case "KeyboardInterrupt", "SystemExit", "GeneratorExit":
		return true
	default:
		return false
	}
}
