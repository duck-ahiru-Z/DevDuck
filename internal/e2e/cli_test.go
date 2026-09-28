package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func buildCLI(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "duck")
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	cmd := exec.Command("go", "build", "-o", path, "./cmd/duck")
	cmd.Dir = filepath.Join("..", "..")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, output)
	}
	return path
}

func TestVersionCLI(t *testing.T) {
	path := buildCLI(t)
	cmd := exec.Command(path, "version")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("version failed: %v\n%s", err, output)
	}
}

func TestInvalidArgumentsCLI(t *testing.T) {
	path := buildCLI(t)
	cmd := exec.Command(path)
	if err := cmd.Run(); err == nil {
		t.Fatal("expected invalid arguments to fail")
	}
}

func TestE2EDoesNotUseNetwork(t *testing.T) {
	if os.Getenv("DEVDUCK_E2E_NETWORK") != "" {
		t.Fatal("network-enabled E2E is not allowed")
	}
}
