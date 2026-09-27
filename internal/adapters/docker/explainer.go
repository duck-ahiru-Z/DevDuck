package docker

import (
	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"regexp"
)

type rule struct {
	kind        string
	pattern     *regexp.Regexp
	explanation model.Explanation
}

func exp(summary, first, second string) model.Explanation {
	return model.Explanation{Summary: summary, Hints: []string{first, second}}
}

var rules = []rule{
	{"DaemonError", regexp.MustCompile(`(?i)cannot connect|error during connect`), exp("Docker CLIがDocker Engineへ接続できていません。", "Docker EngineやDocker Desktopが起動しているか確認してみよう。", "現在のDocker contextや接続先を確認してみよう。")},
	{"PermissionError", regexp.MustCompile(`(?i)permission denied`), exp("Docker Engineへの接続権限がありません。", "実行ユーザーとDockerソケットの権限を確認してみよう。", "Dockerのアクセス設定が現在の環境に合っているか考えてみよう。")},
	{"ImageError", regexp.MustCompile(`(?i)unable to find|repository does not exist|manifest unknown`), exp("指定したイメージを取得または見つけられません。", "イメージ名とタグを確認してみよう。", "使用しているレジストリにそのイメージが存在するか考えてみよう。")},
	{"RegistryError", regexp.MustCompile(`(?i)access denied|unauthorized`), exp("レジストリからイメージを取得する権限がありません。", "対象レジストリと認証状態を確認してみよう。", "そのイメージへのアクセス権があるか考えてみよう。")},
	{"ContainerError", regexp.MustCompile(`(?i)no such container|already in use`), exp("指定したコンテナを見つけられないか、名前が競合しています。", "コンテナ名と現在のコンテナ一覧を確認してみよう。", "既存コンテナとの名前の関係を考えてみよう。")},
	{"NetworkError", regexp.MustCompile(`(?i)network|address already in use|allocated`), exp("Dockerのネットワークやポートを確保できませんでした。", "対象のネットワークやホストポートを確認してみよう。", "別のコンテナや設定と競合していないか考えてみよう。")},
	{"BuildError", regexp.MustCompile(`(?i)failed to solve|dockerfile|copy failed|cache key`), exp("Docker buildの入力やBuildKitの処理に失敗しました。", "Dockerfileとbuild contextの範囲を確認してみよう。", "参照しているファイルや前段のbuild結果を考えてみよう。")},
	{"ComposeError", regexp.MustCompile(`(?i)configuration file|undefined service|undefined network|invalid compose`), exp("Compose設定を読み込めないか、サービス間の参照が成立していません。", "Composeファイルとサービス名を確認してみよう。", "ネットワークや依存サービスの定義が一致しているか考えてみよう。")},
}

type Explainer struct{}

func NewExplainer() *Explainer { return &Explainer{} }
func (e *Explainer) Explain(err model.ErrorInfo) (model.Explanation, bool) {
	if err.Source != "docker" {
		return model.Explanation{}, false
	}
	for _, candidate := range rules {
		if candidate.kind == err.Kind && candidate.pattern.MatchString(err.Message) {
			return candidate.explanation, true
		}
	}
	return model.Explanation{}, false
}
