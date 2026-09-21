package python

import "github.com/duck-ahiru-Z/DevDuck/internal/model"

type Explainer struct{}

func NewExplainer() *Explainer {
	return &Explainer{}
}

var explanations = map[string]model.Explanation{
	"ZeroDivisionError": {
		Summary: "0で割ろうとしたため、Pythonが処理を続けられませんでした。",
		Hints: []string{
			"割る側の値が0になっていないか確認してみよう。",
			"その値がどこで代入されているか辿ってみよう。",
		},
	},
	"NameError": {
		Summary: "Pythonが、使おうとした名前を見つけられていません。",
		Hints: []string{
			"変数名や関数名のスペルを確認してみよう。",
			"その名前を使用する前に定義しているか確認してみよう。",
		},
	},

	"TypeError": {
		Summary: "値の型と、その操作が期待している型が合っていない可能性があります。",
		Hints: []string{
			"エラーが発生した行で使っている値の型を確認してみよう。",
			"type() を使うと、実際の型を確認できます。",
		},
	},

	"IndexError": {
		Summary: "リストなどの存在しない位置を参照している可能性があります。",
		Hints: []string{
			"指定しているインデックスを確認してみよう。",
			"len() を使って要素数を確認してみよう。",
		},
	},

	"KeyError": {
		Summary: "辞書に存在しないキーを参照している可能性があります。",
		Hints: []string{
			"指定しているキーが辞書に存在するか確認してみよう。",
			"キー名のスペルも確認してみよう。",
		},
	},

	"AttributeError": {
		Summary: "その値には、使おうとしている属性やメソッドが存在しない可能性があります。",
		Hints: []string{
			"その変数がどの型なのか確認してみよう。",
			"呼び出しているメソッド名のスペルを確認してみよう。",
		},
	},

	"ModuleNotFoundError": {
		Summary: "Pythonが指定されたモジュールを見つけられていません。",
		Hints: []string{
			"モジュール名のスペルを確認してみよう。",
			"現在使用しているPython環境に、そのモジュールが存在するか確認してみよう。",
		},
	},

	"FileNotFoundError": {
		Summary: "Pythonが指定されたファイルを見つけられていません。",
		Hints: []string{
			"ファイル名やパスが正しいか確認してみよう。",
			"現在の作業ディレクトリも確認してみよう。",
		},
	},

	"SyntaxError": {
		Summary: "Pythonの文法として解釈できない部分があります。",
		Hints: []string{
			"エラーが示している行の記号や括弧を確認してみよう。",
			"その直前の行にも文法ミスがないか確認してみよう。",
		},
	},

	"IndentationError": {
		Summary: "インデントの位置や深さに問題があります。",
		Hints: []string{
			"エラー周辺の行のインデントを比較してみよう。",
			"タブとスペースが混ざっていないか確認してみよう。",
		},
	},
}

func (e *Explainer) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	if err.Source != "python" {
		return model.Explanation{}, false
	}

	explanations, ok := explanations[err.Kind]

	if !ok {
		return model.Explanation{}, false
	}

	return explanations, true
}
