package git

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

func TestExplainKnownDiagnostics(t *testing.T) {
	for _, name := range []string{"not_repository", "pathspec", "auth_failed", "publickey", "dns", "timeout", "non_fast_forward", "merge_conflict", "local_changes", "branch_exists", "remote_exists", "refspec", "unrelated", "identity", "dubious"} {
		t.Run(name, func(t *testing.T) {
			data, _ := os.ReadFile(filepath.Join("testdata", name+".txt"))
			info, ok := Parse(string(data))
			if !ok {
				t.Fatal("parse failed")
			}
			result, explained := NewExplainer().Explain(info)
			if !explained || result.Summary == "" || len(result.Hints) == 0 {
				t.Fatalf("result=%#v", result)
			}
		})
	}
}
func TestUnknownDiagnosticFallsBack(t *testing.T) {
	if _, ok := NewExplainer().Explain(model.ErrorInfo{Source: "git", Kind: "GitError", Message: "future diagnostic"}); ok {
		t.Fatal("expected fallback")
	}
}
