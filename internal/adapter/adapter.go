package adapter

import "github.com/duck-ahiru-Z/DevDuck/internal/model"

// Adapter は、DevDuckが各言語やツールを扱うための共通インターフェース。
type Adapter interface {
	// Adapterの名前を返す。
	Name() string

	// 実行されたコマンドを、このAdapterが処理できるか判定する。
	Detect(
		command string,
		args []string,
	) bool

	// stderrを解析して、共通のErrorInfoに変換する。
	Parse(
		stderr string,
	) (model.ErrorInfo, bool)

	// エラーに対するローカル解説を返す。
	Explain(
		err model.ErrorInfo,
	) (model.Explanation, bool)
}
