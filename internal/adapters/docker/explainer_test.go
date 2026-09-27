package docker

import (
	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"os"
	"path/filepath"
	"testing"
)

func TestExplainKnown(t *testing.T) {
	for _, n := range []string{"daemon", "permission", "image", "pull_denied", "manifest", "container", "conflict", "port", "dockerfile", "solve", "copy", "network", "compose", "compose_service", "unauthorized"} {
		t.Run(n, func(t *testing.T) {
			d, _ := os.ReadFile(filepath.Join("testdata", n+".txt"))
			i, ok := Parse(string(d))
			if !ok {
				t.Fatal("parse")
			}
			e, ok := NewExplainer().Explain(i)
			if !ok || e.Summary == "" || len(e.Hints) == 0 {
				t.Fatalf("%#v", e)
			}
		})
	}
}
func TestUnknownFallsBack(t *testing.T) {
	if _, ok := NewExplainer().Explain(model.ErrorInfo{Source: "docker", Kind: "DockerError", Message: "future"}); ok {
		t.Fatal("expected fallback")
	}
}
