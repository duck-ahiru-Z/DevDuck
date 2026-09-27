package python

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFixtures(t *testing.T) {
	cases := []struct {
		name, kind, message, file string
		line                      int
	}{
		{"zerodivisionerror", "ZeroDivisionError", "division by zero", "app.py", 4},
		{"nameerror", "NameError", "name 'total' is not defined", "app.py", 2},
		{"typeerror", "TypeError", `can only concatenate str (not "int") to str`, "app.py", 3},
		{"indexerror", "IndexError", "list index out of range", "app.py", 2},
		{"keyerror", "KeyError", "'name'", "app.py", 2},
		{"attributeerror", "AttributeError", "'tuple' object has no attribute 'append'", "app.py", 2},
		{"modulenotfounderror", "ModuleNotFoundError", "No module named 'missing_package'", "app.py", 1},
		{"filenotfounderror", "FileNotFoundError", "[Errno 2] No such file or directory: 'missing.txt'", "app.py", 2},
		{"valueerror", "ValueError", "invalid literal for int() with base 10: 'abc'", "app.py", 2},
		{"runtimeerror", "RuntimeError", "worker stopped", "app.py", 2},
		{"syntaxerror", "SyntaxError", "'(' was never closed", "test.py", 3},
		{"indentationerror", "IndentationError", "expected an indented block after 'if' statement on line 2", "test.py", 3},
		{"importerror", "ImportError", "cannot import name 'missing' from 'package'", "app.py", 1},
		{"unboundlocalerror", "UnboundLocalError", "cannot access local variable 'value' where it is not associated with a value", "app.py", 3},
		{"recursionerror", "RecursionError", "maximum recursion depth exceeded", "app.py", 2},
		{"permissionerror", "PermissionError", "[Errno 13] Permission denied: 'secret.txt'", "app.py", 2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tc.name+".txt"))
			if err != nil {
				t.Fatal(err)
			}
			info, ok := Parse(string(data))
			if !ok {
				t.Fatal("expected fixture to be parsed")
			}
			if info.Kind != tc.kind || info.Message != tc.message || info.File != tc.file || info.Line != tc.line {
				t.Fatalf("got kind=%q message=%q file=%q line=%d", info.Kind, info.Message, info.File, info.Line)
			}
		})
	}
}

func TestParseChainedExceptionUsesFinalError(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "chained.txt"))
	if err != nil {
		t.Fatal(err)
	}
	info, ok := Parse(string(data))
	if !ok {
		t.Fatal("expected chained fixture to be parsed")
	}
	if info.Kind != "RuntimeError" || info.Message != "loading failed" || info.File != "app.py" || info.Line != 6 {
		t.Fatalf("got %#v", info)
	}
}
