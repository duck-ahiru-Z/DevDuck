package python

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

func TestExplainSupportedFixturesLocally(t *testing.T) {
	kinds := []string{"zerodivisionerror", "nameerror", "typeerror", "indexerror", "keyerror", "attributeerror", "modulenotfounderror", "filenotfounderror", "valueerror", "runtimeerror", "syntaxerror", "indentationerror", "importerror", "unboundlocalerror", "recursionerror", "permissionerror"}
	for _, name := range kinds {
		t.Run(name, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("testdata", name+".txt"))
			if err != nil {
				t.Fatal(err)
			}
			info, ok := Parse(string(data))
			if !ok {
				t.Fatal("expected fixture to be parsed")
			}
			explanation, explained := NewExplainer().Explain(info)
			if !explained || explanation.Summary == "" || len(explanation.Hints) == 0 {
				t.Fatalf("explanation=%#v explained=%v", explanation, explained)
			}
		})
	}
}

func TestExplainUnknownKindFallsBack(t *testing.T) {
	if _, ok := NewExplainer().Explain(model.ErrorInfo{Source: "python", Kind: "CustomError"}); ok {
		t.Fatal("expected unknown kind to require fallback")
	}
}
