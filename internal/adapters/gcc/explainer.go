package gcc

import (
	"regexp"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

type rule struct {
	pattern     *regexp.Regexp
	explanation model.Explanation
}

var rules = []rule{
	{regexp.MustCompile(`expected ';'`), model.Explanation{Summary: "文の終わりに必要なセミコロンを見つけられませんでした。", Hints: []string{"エラー行と、その直前の文末を確認してみよう。", "その文がどこで終わるべきか考えてみよう。"}}},
	{regexp.MustCompile(`undeclared|undeclared identifier`), model.Explanation{Summary: "使われた名前が、その場所で宣言されていません。", Hints: []string{"名前のスペルと宣言場所を確認してみよう。", "その宣言が現在のスコープから見えるか考えてみよう。"}}},
	{regexp.MustCompile(`implicit declaration of function`), model.Explanation{Summary: "関数の宣言を確認できないまま呼び出しています。", Hints: []string{"必要なヘッダを読み込んでいるか確認してみよう。", "呼び出している関数の宣言がどこにあるか探してみよう。"}}},
	{regexp.MustCompile(`too few arguments|too many arguments`), model.Explanation{Summary: "関数に渡した引数の数が宣言と一致していません。", Hints: []string{"呼び出し側の引数の数を確認してみよう。", "関数宣言が要求する引数を確認してみよう。"}}},
	{regexp.MustCompile(`incompatible types`), model.Explanation{Summary: "代入や式で組み合わせた型に互換性がありません。", Hints: []string{"左右の値の型を確認してみよう。", "その処理が期待する型を宣言と照らし合わせてみよう。"}}},
	{regexp.MustCompile(`expected expression`), model.Explanation{Summary: "値や式が必要な位置に、解釈できる式がありません。", Hints: []string{"エラー位置の記号や演算子の前後を確認してみよう。", "その場所でどんな値を計算するはずか考えてみよう。"}}},
	{regexp.MustCompile(`redefinition`), model.Explanation{Summary: "同じ名前を同じスコープで複数回定義しています。", Hints: []string{"同じ名前の宣言を検索してみよう。", "意図した宣言と重複した宣言を区別してみよう。"}}},
	{regexp.MustCompile(`conflicting types`), model.Explanation{Summary: "同じ関数や名前に異なる型情報があります。", Hints: []string{"宣言と定義の型を比較してみよう。", "ヘッダと実装の宣言が一致しているか確認してみよう。"}}},
	{regexp.MustCompile(`undefined reference to`), model.Explanation{Summary: "呼び出した名前の実体をリンク時に見つけられませんでした。", Hints: []string{"関数の定義がビルド対象に含まれているか確認してみよう。", "必要なライブラリやリンク順を確認してみよう。"}}},
	{regexp.MustCompile(`multiple definition of`), model.Explanation{Summary: "同じ名前の実体が複数の場所で定義されています。", Hints: []string{"重複している定義の場所を確認してみよう。", "宣言と実体の役割を分けられるか考えてみよう。"}}},
}

func (e *Explainer) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	if err.Source != "gcc" {
		return model.Explanation{}, false
	}
	for _, candidate := range rules {
		if candidate.pattern.MatchString(err.Message) {
			return candidate.explanation, true
		}
	}
	return model.Explanation{}, false
}

func NewExplainer() *Explainer { return &Explainer{} }

type Explainer struct{}
