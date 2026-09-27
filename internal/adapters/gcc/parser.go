package gcc

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

var diagnosticPattern = regexp.MustCompile(`^(.*?):([0-9]+)(?::([0-9]+))?:\s+(error|warning|note):\s+(.*)$`)
var linkerPattern = regexp.MustCompile(`^(?:.*?:\s*)?(undefined reference to|multiple definition of)\s+(.+)$`)

func Parse(stderr string) (model.ErrorInfo, bool) {
	var fallback model.ErrorInfo
	for _, line := range strings.Split(stderr, "\n") {
		line = strings.TrimSuffix(line, "\r")
		if match := diagnosticPattern.FindStringSubmatch(line); len(match) > 0 {
			if match[4] == "warning" || match[4] == "note" {
				continue
			}
			lineNumber, _ := strconv.Atoi(match[2])
			fallback = model.ErrorInfo{Source: "gcc", Kind: "CompileError", Message: match[5], File: match[1], Line: lineNumber, Raw: stderr}
			return fallback, true
		}
		if match := linkerPattern.FindStringSubmatch(line); len(match) > 0 {
			return model.ErrorInfo{Source: "gcc", Kind: "LinkerError", Message: match[1] + " " + match[2], Raw: stderr}, true
		}
	}
	return fallback, false
}
