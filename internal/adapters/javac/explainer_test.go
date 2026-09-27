package javac

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

func TestExplainKnownDiagnostics(t *testing.T) {
	for _, name := range []string{"missing_semicolon", "cannot_find_variable", "cannot_find_method", "cannot_find_class", "incompatible_types", "wrong_arguments", "duplicate_variable", "public_filename", "package_missing", "cannot_access", "static_context", "unreported_exception", "missing_return", "unexpected_eof"} {
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
	if _, ok := NewExplainer().Explain(model.ErrorInfo{Source: "javac", Kind: "CompileError", Message: "new javac diagnostic"}); ok {
		t.Fatal("expected fallback")
	}
}

func TestCannotFindSymbolUsesPrimaryDetail(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("testdata", "symbol_scope.txt"))
	if err != nil {
		t.Fatal(err)
	}
	info, ok := Parse(string(data))
	if !ok {
		t.Fatal("expected parse")
	}
	explanation, ok := NewExplainer().Explain(info)
	if !ok || !strings.Contains(explanation.Summary, "メソッド") {
		t.Fatalf("explanation=%#v detail=%q", explanation, info.Detail)
	}
}
