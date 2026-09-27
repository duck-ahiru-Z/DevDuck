package git

import (
	"path/filepath"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

type Adapter struct{ explainer *Explainer }

func NewAdapter() *Adapter      { return &Adapter{explainer: NewExplainer()} }
func (a *Adapter) Name() string { return "git" }
func (a *Adapter) Detect(command string, _ []string) bool {
	switch strings.ToLower(filepath.Base(command)) {
	case "git", "git.exe":
		return true
	default:
		return false
	}
}
func (a *Adapter) Parse(stderr string) (model.ErrorInfo, bool) { return Parse(stderr) }
func (a *Adapter) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	return a.explainer.Explain(err)
}
