package git

import (
	"regexp"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
)

type rule struct {
	kind        string
	pattern     *regexp.Regexp
	explanation model.Explanation
}

var rules = []rule{
	{"RepositoryError", regexp.MustCompile(`not a git repository`), explanation("現在の場所はGitリポジトリとして認識されていません。", "作業ディレクトリと.gitの位置を確認してみよう。", "対象のリポジトリで操作しているか考えてみよう.")},
	{"RefError", regexp.MustCompile(`pathspec`), explanation("指定したパスやブランチをGitが見つけられません。", "名前のスペルと現在の参照を確認してみよう。", "その参照が現在の状態に存在するか考えてみよう.")},
	{"AuthenticationError", regexp.MustCompile(`(?i)authentication failed|publickey`), explanation("リモートへの認証に失敗しました。", "使用している認証方法と資格情報を確認してみよう。", "接続先とアカウントの組み合わせが正しいか考えてみよう.")},
	{"NetworkError", regexp.MustCompile(`(?i)resolve host|timed out|failed to connect`), explanation("Gitリモートへネットワーク接続できませんでした。", "ホスト名とネットワーク接続を確認してみよう。", "リモート側が利用可能か、別の原因と切り分けてみよう.")},
	{"RefError", regexp.MustCompile(`non-fast-forward|remote contains work|\[rejected\]`), explanation("リモート側にローカルが持っていない変更があり、そのまま更新できません。", "ローカルとリモートの履歴差分を確認してみよう。", "先にリモートの変更を取り込む必要があるか考えてみよう.")},
	{"MergeConflict", regexp.MustCompile(`CONFLICT|merge conflict`), explanation("複数の変更が同じ箇所にあり、自動で統合できませんでした。", "競合しているファイルと箇所を確認してみよう。", "各変更をどう組み合わせるか自分の意図を整理してみよう.")},
	{"WorktreeError", regexp.MustCompile(`local changes.*overwritten`), explanation("ローカルの変更が別の操作で上書きされる可能性があります。", "未コミットの変更と対象ファイルを確認してみよう。", "変更をどう扱うか決めてから次の操作を考えてみよう.")},
	{"RefError", regexp.MustCompile(`branch .* already exists|remote .* already exists|refspec .* does not match`), explanation("指定したGitの参照を作成または解決できません。", "ブランチ名やremote名、対象の参照を確認してみよう。", "現在の参照一覧と指定した名前の関係を考えてみよう.")},
	{"RefError", regexp.MustCompile(`unrelated histories`), explanation("2つの履歴に共通の祖先がなく、そのまま統合できません。", "統合しようとしているリポジトリの履歴を確認してみよう。", "別々に始まった履歴をどう扱うか整理してみよう.")},
	{"ConfigError", regexp.MustCompile(`(?i)please tell me who you are|user\.name|user\.email|author identity unknown`), explanation("コミット作成者のIdentityが設定されていません。", "現在のGit設定で名前とメールアドレスを確認してみよう。", "このリポジトリで使うIdentityの範囲を考えてみよう.")},
	{"RepositoryError", regexp.MustCompile(`dubious ownership`), explanation("Gitがリポジトリの所有者を安全だと確認できません。", "リポジトリの場所と実行ユーザーを確認してみよう。", "安全性を確認できる環境で操作しているか考えてみよう.")},
}

func explanation(summary, first, second string) model.Explanation {
	return model.Explanation{Summary: summary, Hints: []string{first, second}}
}
func NewExplainer() *Explainer { return &Explainer{} }

type Explainer struct{}

func (e *Explainer) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	if err.Source != "git" {
		return model.Explanation{}, false
	}
	for _, candidate := range rules {
		if candidate.kind == err.Kind && candidate.pattern.MatchString(err.Raw) {
			return candidate.explanation, true
		}
	}
	return model.Explanation{}, false
}
