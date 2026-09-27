package javac

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

var diagnosticPattern = regexp.MustCompile(`^(.*\.java):([0-9]+):\s+error:\s+(.*)$`)

func Parse(stderr string) (model.ErrorInfo, bool) {
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSuffix(line, "\r")
		match := diagnosticPattern.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		lineNumber, _ := strconv.Atoi(match[2])
		return model.ErrorInfo{Source: "javac", Kind: "CompileError", Message: match[3], File: match[1], Line: lineNumber, Raw: stderr}, true
	}
	return model.ErrorInfo{}, false
}
