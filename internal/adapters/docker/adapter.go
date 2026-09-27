package docker

import (
	"path/filepath"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

type Adapter struct{ explainer *Explainer }

func NewAdapter() *Adapter      { return &Adapter{explainer: NewExplainer()} }
func (a *Adapter) Name() string { return "docker" }
func (a *Adapter) Detect(command string, _ []string) bool {
	switch strings.ToLower(filepath.Base(command)) {
	case "docker", "docker.exe", "docker-compose", "docker-compose.exe":
		return true
	default:
		return false
	}
}
func (a *Adapter) Parse(stderr string) (model.ErrorInfo, bool) { return Parse(stderr) }
func (a *Adapter) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	return a.explainer.Explain(err)
}
