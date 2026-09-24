package python

import (
	"path/filepath"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

// Adapter はPython用のAdapter。
type Adapter struct {
	explainer *Explainer
}

// NewAdapter はPython Adapterを作成する。
func NewAdapter() *Adapter {
	return &Adapter{
		explainer: NewExplainer(),
	}
}

// Name はAdapter名を返す。
func (a *Adapter) Name() string {
	return "python"
}

// Detect はPython系のコマンドかどうかを判定する。
func (a *Adapter) Detect(
	command string,
	args []string,
) bool {
	name := strings.ToLower(
		filepath.Base(command),
	)

	switch name {
	case "python",
		"python.exe",
		"python3",
		"python3.exe",
		"py",
		"py.exe":
		return true
	}

	return false
}

// Parse はPythonのstderrを解析する。
func (a *Adapter) Parse(
	stderr string,
) (model.ErrorInfo, bool) {
	return Parse(stderr)
}

// Explain はPythonエラーのローカル解説を返す。
func (a *Adapter) Explain(
	err model.ErrorInfo,
) (model.Explanation, bool) {
	return a.explainer.Explain(err)
}
