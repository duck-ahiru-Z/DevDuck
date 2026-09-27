package javac

import (
	"regexp"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

type rule struct {
	pattern     *regexp.Regexp
	explanation model.Explanation
}

var rules = []rule{
	{regexp.MustCompile(`^';' expected$`), model.Explanation{Summary: "文の区切りに必要なセミコロンを見つけられませんでした。", Hints: []string{"エラー行と直前の文末を確認してみよう。", "その文がどこで終わるべきか考えてみよう。"}}},
	{regexp.MustCompile(`^incompatible types:`), model.Explanation{Summary: "代入や式で、互換性のない型を組み合わせています。", Hints: []string{"代入元と代入先の型を確認してみよう。", "その値を扱う設計上の型を考えてみよう。"}}},
	{regexp.MustCompile(`method .* cannot be applied to given types`), model.Explanation{Summary: "メソッド呼び出しの引数が宣言と一致していません。", Hints: []string{"requiredとfoundの引数を比較してみよう。", "呼び出し先のメソッド宣言を確認してみよう。"}}},
	{regexp.MustCompile(`variable .* is already defined`), model.Explanation{Summary: "同じスコープで変数を重複して定義しています。", Hints: []string{"同じ名前の宣言を検索してみよう。", "変数の有効範囲を確認してみよう。"}}},
	{regexp.MustCompile(`class .* is public, should be declared in a file named`), model.Explanation{Summary: "publicクラス名とJavaファイル名が一致していません。", Hints: []string{"クラス宣言とファイル名を比較してみよう。", "Javaのpublicクラスの命名規則を確認してみよう。"}}},
	{regexp.MustCompile(`^package .* does not exist$`), model.Explanation{Summary: "指定したパッケージをコンパイラが見つけられていません。", Hints: []string{"importの名前と依存ライブラリを確認してみよう。", "classpathに必要なライブラリが含まれているか考えてみよう。"}}},
	{regexp.MustCompile(`^cannot access`), model.Explanation{Summary: "参照先のクラスへアクセスできません。", Hints: []string{"クラスファイルや依存関係が正しく配置されているか確認してみよう。", "アクセス修飾子とclasspathを確認してみよう。"}}},
	{regexp.MustCompile(`non-static .* referenced from a static context`), model.Explanation{Summary: "staticな場所からインスタンス用のメンバーを参照しています。", Hints: []string{"参照元がstaticか確認してみよう。", "対象メンバーがどのインスタンスに属するか考えてみよう。"}}},
	{regexp.MustCompile(`^unreported exception .* must be caught or declared`), model.Explanation{Summary: "チェック例外を処理せずに呼び出しています。", Hints: []string{"例外を処理する呼び出し箇所を確認してみよう。", "メソッド宣言で例外を伝播させる設計か考えてみよう。"}}},
	{regexp.MustCompile(`^missing return statement$`), model.Explanation{Summary: "値を返すメソッドの一部の経路にreturnがありません。", Hints: []string{"すべての分岐で戻り値が決まるか確認してみよう。", "メソッドの戻り値型と条件分岐を見比べてみよう。"}}},
	{regexp.MustCompile(`^reached end of file while parsing$`), model.Explanation{Summary: "ファイル末尾までに閉じられていない構文があります。", Hints: []string{"括弧や波括弧の対応を確認してみよう。", "エラー行より前のブロック開始位置を確認してみよう。"}}},
}

func (e *Explainer) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	if err.Source != "javac" || err.Kind != "CompileError" {
		return model.Explanation{}, false
	}
	if err.Message == "cannot find symbol" {
		if strings.Contains(err.Detail, "symbol:   variable") {
			return explanation("変数が見つかりません。", "変数名のスペルと宣言場所を確認してみよう。", "その宣言が現在のスコープから見えるか考えてみよう."), true
		}
		if strings.Contains(err.Detail, "symbol:   method") {
			return explanation("メソッドが見つかりません。", "メソッド名と引数を確認してみよう。", "対象の型にそのメソッドが定義されているか考えてみよう."), true
		}
		if strings.Contains(err.Detail, "symbol:   class") {
			return explanation("クラスが見つかりません。", "クラス名とimportを確認してみよう。", "必要な依存関係がclasspathにあるか考えてみよう."), true
		}
	}
	for _, candidate := range rules {
		if candidate.pattern.MatchString(err.Message) {
			return candidate.explanation, true
		}
	}
	return model.Explanation{}, false
}

func explanation(summary, first, second string) model.Explanation {
	return model.Explanation{Summary: summary, Hints: []string{first, second}}
}
func NewExplainer() *Explainer { return &Explainer{} }

type Explainer struct{}
