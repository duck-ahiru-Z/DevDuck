package javac

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	adapter := NewAdapter()
	for _, command := range []string{"javac", "javac.exe", `/usr/bin/javac`} {
		if !adapter.Detect(command, nil) {
			t.Errorf("Detect(%q) = false", command)
		}
	}
	if adapter.Detect("java", nil) {
		t.Error("java runtime should not be handled")
	}
}

func TestParseFixtures(t *testing.T) {
	cases := []struct {
		name, message, file string
		line                int
	}{
		{"missing_semicolon", "';' expected", "Main.java", 5},
		{"cannot_find_variable", "cannot find symbol", "Main.java", 7},
		{"cannot_find_method", "cannot find symbol", "Main.java", 9},
		{"cannot_find_class", "cannot find symbol", "Main.java", 2},
		{"incompatible_types", "incompatible types: String cannot be converted to int", "Main.java", 4},
		{"wrong_arguments", "method add in class Calc cannot be applied to given types;", "Main.java", 8},
		{"duplicate_variable", "variable count is already defined in method run()", "Main.java", 6},
		{"public_filename", "class PublicApp is public, should be declared in a file named PublicApp.java", "Main.java", 1},
		{"package_missing", "package missing.lib does not exist", "Main.java", 1},
		{"cannot_access", "cannot access Secret", "Main.java", 3},
		{"static_context", "non-static method run() cannot be referenced from a static context", "Main.java", 5},
		{"unreported_exception", "unreported exception IOException; must be caught or declared to be thrown", "Main.java", 7},
		{"missing_return", "missing return statement", "Main.java", 4},
		{"unexpected_eof", "reached end of file while parsing", "Main.java", 10},
		{"multiple_errors", "';' expected", "Main.java", 3},
		{"windows_path", "';' expected", `C:\Users\name\project\Main.java`, 5},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", tc.name+".txt"))
			if err != nil {
				t.Fatal(err)
			}
			info, ok := Parse(string(data))
			if !ok || info.Kind != "CompileError" || info.Message != tc.message || info.File != tc.file || info.Line != tc.line {
				t.Fatalf("got %#v, ok=%v", info, ok)
			}
		})
	}
}
