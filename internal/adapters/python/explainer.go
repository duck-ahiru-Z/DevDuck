package python

import (
	"regexp"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

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
	"ImportError": {
		Summary: "モジュールから指定した名前を読み込めませんでした。",
		Hints:   []string{"名前のスペルと、モジュールが公開している名前を確認してみよう。", "実行しているPython環境を確認してみよう。"},
	},
	"RuntimeError": {
		Summary: "実行中に一般的な問題が発生しました。",
		Hints:   []string{"エラーメッセージと直前の処理を確認してみよう。", "どの条件でこの処理に到達するか考えてみよう。"},
	},
	"ValueError": {
		Summary: "値そのものが、処理で期待される形式や範囲に合っていません。",
		Hints:   []string{"入力値や変換対象の内容を確認してみよう。", "どの値がこの処理に渡っているか確認してみよう。"},
	},
	"UnboundLocalError": {
		Summary: "関数内のローカル変数を、値が設定される前に使おうとしています。",
		Hints:   []string{"その変数がすべての分岐で代入されるか確認してみよう。", "関数の外側の変数との名前の重なりも確認してみよう。"},
	},
	"RecursionError": {
		Summary: "関数の呼び出しが深くなりすぎました。",
		Hints:   []string{"再帰呼び出しが終了条件に到達するか確認してみよう。", "同じ引数で呼び出し続けていないか確認してみよう。"},
	},
	"PermissionError": {
		Summary: "ファイルやリソースへアクセスする権限がありません。",
		Hints:   []string{"対象のパスとアクセスモードを確認してみよう。", "実行ユーザーに必要な権限があるか確認してみよう。"},
	},
}

type messageRule struct {
	kind    string
	pattern *regexp.Regexp
	result  model.Explanation
}

var messageRules = []messageRule{
	{kind: "TypeError", pattern: regexp.MustCompile(`unsupported operand type`), result: model.Explanation{Summary: "演算子が、その値の型の組み合わせに対応していません。", Hints: []string{"演算子の左右の値の型を確認してみよう。", "暗黙の型変換を期待していないか考えてみよう。"}}},
	{kind: "TypeError", pattern: regexp.MustCompile(`not subscriptable`), result: model.Explanation{Summary: "添字で要素を取り出せない型を、添字付きで使おうとしています。", Hints: []string{"添字を使っている値の型を確認してみよう。", "その値が本当にシーケンスや辞書か考えてみよう。"}}},
	{kind: "KeyError", pattern: regexp.MustCompile(`^'[^']+'$`), result: model.Explanation{Summary: "辞書に指定したキーが存在しません。", Hints: []string{"実際のキー一覧と指定値を比較してみよう。", "キーが作られる条件も確認してみよう。"}}},
}

func (e *Explainer) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	if err.Source != "python" {
		return model.Explanation{}, false
	}

	for _, rule := range messageRules {
		if rule.kind == err.Kind && rule.pattern.MatchString(err.Message) {
			return rule.result, true
		}
	}

	explanation, ok := explanations[err.Kind]

	if !ok {
		return model.Explanation{}, false
	}

	return explanation, true
}
