package doctor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/duck-ahiru-Z/DevDuck/internal/buildinfo"
	"github.com/duck-ahiru-Z/DevDuck/internal/config"
	"github.com/duck-ahiru-Z/DevDuck/internal/credential"
)

type Checker struct {
	Credential credential.Store
	LookPath   func(string) (string, error)
	ConfigLoad func() error
}

func NewChecker() Checker {
	return Checker{Credential: credential.NewKeyringStore(), LookPath: exec.LookPath, ConfigLoad: func() error { _, err := config.Load(); return err }}
}

func Run(out io.Writer) int { return NewChecker().Run(out) }

func (c Checker) Run(out io.Writer) int {
	if c.LookPath == nil {
		c.LookPath = exec.LookPath
	}
	if c.ConfigLoad == nil {
		c.ConfigLoad = func() error { _, err := config.Load(); return err }
	}
	fmt.Fprintln(out, "DevDuck Doctor")
	fmt.Fprintln(out, "")
	fmt.Fprintf(out, "Version: %s\n", buildinfo.Version)
	fmt.Fprintf(out, "OS: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
	configPath, configErr := config.Path()
	cachePath, cacheErr := config.CachePath()
	fmt.Fprintf(out, "Config path: %s\n", configPath)
	fmt.Fprintf(out, "Cache path: %s\n\n", cachePath)
	status(out, "Config", configErr == nil && c.ConfigLoad() == nil)
	status(out, "Cache", cacheErr == nil && cacheAvailable(cachePath))
	credentialOK := false
	if c.Credential != nil {
		_, err := c.Credential.Get("DevDuck", "gemini")
		credentialOK = err == nil || err == credential.ErrNotFound
	}
	status(out, "Credential Store", credentialOK)
	if c.Credential != nil {
		_, err := c.Credential.Get("DevDuck", "gemini")
		if err == nil {
			status(out, "Gemini credential configured", true)
		} else {
			warn(out, "Gemini credential not configured")
		}
	}
	if os.Getenv("DEVDUCK_AI_PROVIDER") == "" || os.Getenv("DEVDUCK_AI_PROVIDER") == "gemini" {
		status(out, "AI provider: Gemini", true)
	} else {
		warn(out, "AI provider: "+os.Getenv("DEVDUCK_AI_PROVIDER"))
	}
	for _, command := range []string{"python", "python3", "py", "gcc", "javac", "git", "docker"} {
		if _, err := c.LookPath(command); err == nil {
			status(out, command, true)
		} else {
			warn(out, command+" not found")
		}
	}
	return 0
}

func cacheAvailable(path string) bool {
	if info, err := os.Stat(filepath.Dir(path)); err == nil {
		return info.IsDir()
	}
	return os.MkdirAll(filepath.Dir(path), 0700) == nil
}
func status(out io.Writer, name string, ok bool) {
	if ok {
		fmt.Fprintf(out, "[OK] %s\n", name)
	} else {
		warn(out, name)
	}
}
func warn(out io.Writer, message string) { fmt.Fprintf(out, "[WARN] %s\n", message) }
