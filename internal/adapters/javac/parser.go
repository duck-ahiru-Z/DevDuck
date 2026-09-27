package javac

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

var diagnosticPattern = regexp.MustCompile(`^(.*\.java):([0-9]+):\s+error:\s+(.*)$`)

func Parse(stderr string) (model.ErrorInfo, bool) {
	lines := strings.Split(stderr, "\n")
	for index, line := range lines {
		line = strings.TrimSuffix(line, "\r")
		match := diagnosticPattern.FindStringSubmatch(line)
		if len(match) == 0 {
			continue
		}
		lineNumber, _ := strconv.Atoi(match[2])
		end := index + 1
		for end < len(lines) && !diagnosticPattern.MatchString(strings.TrimSuffix(lines[end], "\r")) {
			end++
		}
		return model.ErrorInfo{Source: "javac", Kind: "CompileError", Message: match[3], File: match[1], Line: lineNumber, Raw: stderr, Detail: strings.Join(lines[index+1:end], "\n")}, true
	}
	return model.ErrorInfo{}, false
}
