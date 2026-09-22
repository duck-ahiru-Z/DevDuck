package teaching

import (
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

func TestSessionReturnsHintsInOrder(t *testing.T) {
	explanation := model.Explanation{
		Summary: "test",
		Hints: []string{
			"first hint",
			"second hint",
		},
	}

	session := NewSession(explanation)

	first, ok := session.NextHint()

	if !ok {
		t.Fatal("expected first hint")
	}

	if first != "first hint" {
		t.Errorf("expected first hint, got %q", first)
	}

	second, ok := session.NextHint()

	if !ok {
		t.Fatal("expected second hint")
	}

	if second != "second hint" {
		t.Errorf("expected second hint, got %q", second)
	}

	if session.HasNextHint() {
		t.Error("expected no hints remaining")
	}
}
