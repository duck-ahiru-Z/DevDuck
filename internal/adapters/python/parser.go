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
	`(?m)^([A-Za-z_][A-Za-z0-9_]*(?:Error|Exception))(?::\s*(.*))?$`,
)

func Parse(stderr string) (model.ErrorInfo, bool) {
	errorMatches := errorPattern.FindAllStringSubmatch(stderr, -1)

	if len(errorMatches) == 0 {
		return model.ErrorInfo{}, false
	}

	lastError := errorMatches[len(errorMatches)-1]

	info := model.ErrorInfo{
		Source: "python",
		Kind:   lastError[1],
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
