package teaching

import "testing"

func TestPolicyForBeginner(t *testing.T) {
	policy := PolicyForLevel(Beginner)

	if policy.Level != Beginner {
		t.Errorf(
			"expected beginner, got %q",
			policy.Level,
		)
	}

	if !policy.ShowSummary {
		t.Error(
			"expected beginner to show summary",
		)
	}

	if policy.MaxHints != 3 {
		t.Errorf(
			"expected 3 hints, got %d",
			policy.MaxHints,
		)
	}
}

func TestPolicyForAdvanced(t *testing.T) {
	policy := PolicyForLevel(Advanced)

	if policy.ShowSummary {
		t.Error(
			"expected advanced to hide summary",
		)
	}

	if policy.MaxHints != 1 {
		t.Errorf(
			"expected 1 hint, got %d",
			policy.MaxHints,
		)
	}
}
