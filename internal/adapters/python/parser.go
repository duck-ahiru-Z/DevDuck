package python

import (
	"regexp"
	"strconv"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

var filePattern = regexp.MustCompile(
	`File "([^"]+)", line ([0-9]+)`,
)

var errorPattern = regexp.MustCompile(
	`(?m)^((?:[A-Za-z_][A-Za-z0-9_]*\.)*[A-Za-z_][A-Za-z0-9_]*)(?::\s*(.*))$`,
)

var controlPattern = regexp.MustCompile(`(?m)^(KeyboardInterrupt|SystemExit|GeneratorExit)$`)

func Parse(stderr string) (model.ErrorInfo, bool) {
	errorMatches := errorPattern.FindAllStringSubmatch(stderr, -1)
	if len(errorMatches) == 0 {
		controlMatches := controlPattern.FindAllStringSubmatch(stderr, -1)
		if len(controlMatches) > 0 {
			return model.ErrorInfo{Source: "python", Kind: controlMatches[len(controlMatches)-1][1], Raw: stderr}, true
		}
	}

	if len(errorMatches) == 0 {
		return model.ErrorInfo{}, false
	}

	lastError := errorMatches[len(errorMatches)-1]

	info := model.ErrorInfo{
		Source: "python",
		Kind:   unqualifiedKind(lastError[1]),
		Raw:    stderr,
	}

	if len(lastError) >= 3 {
		info.Message = lastError[2]
	}

	fileMatches := filePattern.FindAllStringSubmatch(stderr, -1)

	if len(fileMatches) > 0 {
		lastFile := fileMatches[len(fileMatches)-1]

		info.File = lastFile[1]

		line, err := strconv.Atoi(lastFile[2])

		if err == nil {
			info.Line = line
		}
	}

	return info, true
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
