package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

type result struct {
	code   int
	output string
}

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

func isolatedEnv(t *testing.T) []string {
	t.Helper()
	root := t.TempDir()
	env := os.Environ()
	if runtime.GOOS == "windows" {
		env = append(env, "APPDATA="+filepath.Join(root, "AppData"), "LOCALAPPDATA="+filepath.Join(root, "LocalAppData"), "USERPROFILE="+root)
	} else {
		env = append(env, "HOME="+root, "XDG_CONFIG_HOME="+filepath.Join(root, "config"), "XDG_CACHE_HOME="+filepath.Join(root, "cache"))
	}
	env = append(env, "GEMINI_API_KEY=", "DEVDUCK_AI_PROVIDER=gemini", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL="+filepath.Join(root, "no-global-gitconfig"))
	return env
}

func runCLI(t *testing.T, binary string, env []string, dir string, args ...string) result {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = dir
	cmd.Env = env
	output, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			code = 1
		}
	}
	return result{code: code, output: string(output)}
}

func tool(t *testing.T, names ...string) string {
	t.Helper()
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			if probe := exec.Command(path, "--version"); probe.Run() == nil {
				return path
			}
		}
	}
	t.Skipf("tool not found: %s", strings.Join(names, ", "))
	return ""
}

func TestVersionAndInvalidArgsCLI(t *testing.T) {
	binary := buildCLI(t)
	env := isolatedEnv(t)
	for _, args := range [][]string{{"version"}, {"--version"}} {
		got := runCLI(t, binary, env, t.TempDir(), args...)
		if got.code != 0 || !strings.Contains(got.output, "DevDuck") {
			t.Fatalf("args=%v result=%#v", args, got)
		}
	}
	got := runCLI(t, binary, env, t.TempDir())
	if got.code == 0 {
		t.Fatalf("expected invalid args: %#v", got)
	}
}

func TestConfigAndDoctorCLI(t *testing.T) {
	binary := buildCLI(t)
	env := isolatedEnv(t)
	dir := t.TempDir()
	for _, level := range []string{"beginner", "intermediate", "advanced"} {
		got := runCLI(t, binary, env, dir, "config", "set", "level", level)
		if got.code != 0 {
			t.Fatalf("set %s: %#v", level, got)
		}
		got = runCLI(t, binary, env, dir, "config", "show")
		if got.code != 0 || !strings.Contains(got.output, level) {
			t.Fatalf("show %s: %#v", level, got)
		}
	}
	if got := runCLI(t, binary, env, dir, "config", "set", "level", "invalid"); got.code == 0 {
		t.Fatalf("expected invalid level: %#v", got)
	}
	got := runCLI(t, binary, env, dir, "doctor")
	if got.code != 0 || !strings.Contains(got.output, "DevDuck Doctor") || !strings.Contains(got.output, "Config path") || !strings.Contains(got.output, "Cache path") {
		t.Fatalf("doctor: %#v", got)
	}
}

func TestPythonCLI(t *testing.T) {
	python := tool(t, "python", "python3", "py")
	binary := buildCLI(t)
	env := isolatedEnv(t)
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "success.py"), []byte("print(\"hello\")\n"), 0600)
	got := runCLI(t, binary, env, dir, python, "success.py")
	if got.code != 0 || !strings.Contains(got.output, "hello") || strings.Contains(got.output, "DevDuck detected") {
		t.Fatalf("success: %#v", got)
	}
	os.WriteFile(filepath.Join(dir, "zero.py"), []byte("x = 1 / 0\n"), 0600)
	got = runCLI(t, binary, env, dir, python, "zero.py")
	if got.code == 0 || !strings.Contains(got.output, "ZeroDivisionError") || !strings.Contains(got.output, "Explanation") {
		t.Fatalf("zero: %#v", got)
	}
	os.WriteFile(filepath.Join(dir, "name.py"), []byte("print(missing_name)\n"), 0600)
	got = runCLI(t, binary, env, dir, python, "name.py")
	if got.code == 0 || !strings.Contains(got.output, "NameError") || !strings.Contains(got.output, "Explanation") {
		t.Fatalf("name: %#v", got)
	}
}

func TestGCCCLI(t *testing.T) {
	gcc := tool(t, "gcc")
	binary := buildCLI(t)
	env := isolatedEnv(t)
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "good.c"), []byte("int main(void) { return 0; }\n"), 0600)
	got := runCLI(t, binary, env, dir, gcc, "good.c", "-o", filepath.Join(dir, "good"))
	if got.code != 0 {
		t.Fatalf("gcc success: %#v", got)
	}
	os.WriteFile(filepath.Join(dir, "bad.c"), []byte("int main(void) { return 0 }\n"), 0600)
	got = runCLI(t, binary, env, dir, gcc, "bad.c", "-o", filepath.Join(dir, "bad"))
	if got.code == 0 || !strings.Contains(got.output, "CompileError") {
		t.Fatalf("gcc error: %#v", got)
	}
}

func TestJavacCLI(t *testing.T) {
	javac := tool(t, "javac")
	binary := buildCLI(t)
	env := isolatedEnv(t)
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Good.java"), []byte("class Good { public static void main(String[] a) {} }\n"), 0600)
	got := runCLI(t, binary, env, dir, javac, "Good.java")
	if got.code != 0 {
		t.Fatalf("javac success: %#v", got)
	}
	os.WriteFile(filepath.Join(dir, "Bad.java"), []byte("class Bad { void run() { System.out.println(missing); } }\n"), 0600)
	got = runCLI(t, binary, env, dir, javac, "Bad.java")
	if got.code == 0 || !strings.Contains(got.output, "CompileError") || !strings.Contains(got.output, "Bad.java") {
		t.Fatalf("javac error: %#v", got)
	}
}

func TestGitCLI(t *testing.T) {
	git := tool(t, "git")
	binary := buildCLI(t)
	env := isolatedEnv(t)
	root := t.TempDir()
	got := runCLI(t, binary, env, root, git, "status")
	if got.code != 128 || !strings.Contains(got.output, "RepositoryError") {
		t.Fatalf("not repo: %#v", got)
	}
	repo := filepath.Join(root, "repo")
	os.Mkdir(repo, 0700)
	init := exec.Command(git, "init", repo)
	init.Env = env
	if out, err := init.CombinedOutput(); err != nil {
		t.Fatalf("git init: %v %s", err, out)
	}
	os.WriteFile(filepath.Join(repo, "file.txt"), []byte("x\n"), 0600)
	add := exec.Command(git, "-C", repo, "add", "file.txt")
	add.Env = env
	if out, err := add.CombinedOutput(); err != nil {
		t.Fatalf("git add: %v %s", err, out)
	}
	got = runCLI(t, binary, env, repo, git, "checkout", "definitely-not-existing-branch")
	if got.code == 0 || !strings.Contains(got.output, "RefError") {
		t.Fatalf("pathspec: %#v", got)
	}
	got = runCLI(t, binary, env, repo, git, "commit", "-m", "test")
	if got.code == 0 || (!strings.Contains(got.output, "identity") && !strings.Contains(got.output, "Identity")) {
		t.Fatalf("identity: %#v", got)
	}
}
