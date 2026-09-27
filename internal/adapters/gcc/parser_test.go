package gcc

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	adapter := NewAdapter()
	for _, command := range []string{"gcc", "gcc.exe", "cc", "cc.exe", `/usr/bin/gcc`} {
		if !adapter.Detect(command, nil) {
			t.Errorf("Detect(%q) = false", command)
		}
	}
	if adapter.Detect("clang", nil) {
		t.Error("clang should not be handled by GCC adapter")
	}
}

func TestParseFixtures(t *testing.T) {
	cases := []struct {
		name, kind, message, file string
		line                      int
	}{
		{"missing_semicolon", "CompileError", "expected ';' before 'return'", "main.c", 5},
		{"undeclared", "CompileError", "'count' undeclared (first use in this function)", "main.c", 8},
		{"implicit_function", "CompileError", "implicit declaration of function 'missing'", "main.c", 10},
		{"arguments", "CompileError", "too few arguments to function 'add'", "main.c", 4},
		{"incompatible_types", "CompileError", "incompatible types when assigning to type 'int' from type 'char *'", "main.c", 6},
		{"expected_expression", "CompileError", "expected expression before ';' token", "main.c", 3},
		{"redefinition", "CompileError", "redefinition of 'main'", "main.c", 9},
		{"conflicting_types", "CompileError", "conflicting types for 'run'", "main.c", 7},
		{"undefined_reference", "LinkerError", "undefined reference to `foo'", "", 0},
		{"multiple_definition", "LinkerError", "multiple definition of `main'; main.o:main.c:2: first defined here", "", 0},
		{"multiple_errors", "CompileError", "expected ';' before 'return'", "main.c", 3},
		{"windows_path", "CompileError", "expected ';' before 'return'", `C:\Users\name\project\main.c`, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tc.name+".txt"))
			if err != nil {
				t.Fatal(err)
			}
			info, ok := Parse(string(data))
			if !ok || info.Kind != tc.kind || info.Message != tc.message || info.File != tc.file || info.Line != tc.line {
				t.Fatalf("got %#v, ok=%v", info, ok)
			}
		})
	}
}
