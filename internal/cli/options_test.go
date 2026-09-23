package cli

import (
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

func TestParseWithLevel(t *testing.T) {
	options, err := Parse([]string{
		"--level",
		"beginner",
		"python",
		"test.py",
	})

	if err != nil {
		t.Fatal(err)
	}

	if options.Level != teaching.Beginner {
		t.Errorf(
			"expected beginner, got %q",
			options.Level,
		)
	}

	if options.Command != "python" {
		t.Errorf(
			"expected python, got %q",
			options.Command,
		)
	}

	if len(options.Args) != 1 {
		t.Fatalf(
			"expected 1 argument, got %d",
			len(options.Args),
		)
	}

	if options.Args[0] != "test.py" {
		t.Errorf(
			"expected test.py, got %q",
			options.Args[0],
		)
	}
}
