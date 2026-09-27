package gcc

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

func TestExplainKnownDiagnostics(t *testing.T) {
	for _, name := range []string{"missing_semicolon", "undeclared", "implicit_function", "arguments", "incompatible_types", "expected_expression", "redefinition", "conflicting_types", "undefined_reference", "multiple_definition"} {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", name+".txt"))
			if err != nil {
				t.Fatal(err)
			}
			info, ok := Parse(string(data))
			if !ok {
				t.Fatal("expected parse")
			}
			result, explained := NewExplainer().Explain(info)
			if !explained || result.Summary == "" || len(result.Hints) == 0 {
				t.Fatalf("result=%#v explained=%v", result, explained)
			}
		})
	}
}

func TestUnknownDiagnosticFallsBack(t *testing.T) {
	if _, ok := NewExplainer().Explain(model.ErrorInfo{Source: "gcc", Kind: "CompileError", Message: "some new GCC diagnostic"}); ok {
		t.Fatal("unknown diagnostic should use AI fallback")
	}
}
