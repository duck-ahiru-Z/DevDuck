package runner

import (
	"os"
	"testing"
)

func TestRunPreservesChildExitCodes(t *testing.T) {
	t.Setenv("GO_WANT_HELPER_PROCESS", "1")
	for _, code := range []int{0, 1, 2} {
		result := Run(os.Args[0], []string{"-test.run=TestRunnerHelperProcess", "--", string(rune('0' + code))}, nil, os.Stdout, os.Stderr)
		if result.ExitCode != code {
			t.Fatalf("code=%d got=%d", code, result.ExitCode)
		}
	}
}

func TestRunCommandFailureIsNonZero(t *testing.T) {
	result := Run("devduck-command-that-does-not-exist", nil, nil, os.Stdout, os.Stderr)
	if result.ExitCode == 0 || result.Err == nil {
		t.Fatalf("result=%#v", result)
	}
}

func TestRunnerHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(int(os.Args[len(os.Args)-1][0] - '0'))
}
