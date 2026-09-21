package python

import (
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

func TestExplainZeroDivisionError(t *testing.T) {
	explainer := NewExplainer()

	errInfo := model.ErrorInfo{
		Source: "python",
		Kind:   "ZeroDivisionError",
	}

	explanation, ok := explainer.Explain(errInfo)

	if !ok {
		t.Fatal("expected ZeroDivisionError to be explained")
	}

	if explanation.Summary == "" {
		t.Error("expected Summary to not be empty")
	}

	if len(explanation.Hints) == 0 {
		t.Error("expected at least one hint")
	}
}
