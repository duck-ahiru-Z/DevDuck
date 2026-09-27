package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/adapter"
	"github.com/duck-ahiru-Z/DevDuck/internal/adapters/python"
	duckai "github.com/duck-ahiru-Z/DevDuck/internal/ai"
	"github.com/duck-ahiru-Z/DevDuck/internal/ai/factory"
	"github.com/duck-ahiru-Z/DevDuck/internal/cache"
	"github.com/duck-ahiru-Z/DevDuck/internal/cli"
	"github.com/duck-ahiru-Z/DevDuck/internal/config"
	"github.com/duck-ahiru-Z/DevDuck/internal/credential"
	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/redact"
	"github.com/duck-ahiru-Z/DevDuck/internal/runner"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
	"github.com/duck-ahiru-Z/DevDuck/internal/ui/terminal"
	"golang.org/x/term"
)

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) >= 2 && os.Args[1] == "auth" {
		return handleAuth(os.Args[2:])
	}
	if len(os.Args) >= 2 && os.Args[1] == "config" {
		return handleConfig(os.Args[2:])
	}

	cfg, err := config.Load()

	if err != nil {
		fmt.Println("設定ファイルの読み込みに失敗しました:", err)
		return 1
	}

	options, err := cli.Parse(os.Args[1:])

	if err != nil {
		fmt.Println("DevDuck:", err)
		fmt.Println()
		fmt.Println(
			"使い方: duck [--level beginner|intermediate|advanced] <command> [args...]",
		)
		return 1
	}

	level := cfg.Level

	if options.LevelSet {
		level = options.Level
	}

	command := options.Command
	args := options.Args

	registry := adapter.NewRegistry(
		python.NewAdapter(),
	)

	selectedAdapter, adapterFound := registry.Find(
		command,
		args,
	)

	result := runner.Run(command, args, os.Stdin, os.Stdout, os.Stderr)
	if result.ExitCode == 0 {
		return 0
	}

	fmt.Println()
	fmt.Println("DevDuck detected an error")

	if !adapterFound {
		fmt.Println()
		fmt.Println("このコマンドにはまだ対応していません。")
		return result.ExitCode
	}

	errorInfo, ok := selectedAdapter.Parse(
		result.Stderr,
	)

	if !ok {
		fmt.Println()
		fmt.Println("エラーを解析できませんでした。")
		return result.ExitCode
	}

	fmt.Println()
	fmt.Println("Source :", errorInfo.Source)
	fmt.Println("Type   :", errorInfo.Kind)
	fmt.Println("Message:", errorInfo.Message)
	fmt.Println("File   :", errorInfo.File)
	fmt.Println("Line   :", errorInfo.Line)
	if errorInfo.SkipTeaching {
		return result.ExitCode
	}

	explanation, explained := selectedAdapter.Explain(
		errorInfo,
	)

	// ローカルで説明できない場合だけAIを使う。
	if !explained {
		cachePath, cacheErr := config.CachePath()
		store := (*cache.Store)(nil)
		if cacheErr == nil {
			store = cache.New(cachePath)
			if cached, found, err := store.Get(cache.Key(errorInfo, level, "")); err == nil && found {
				explanation = cached
				explained = true
			}
		}
		if !explained {
			explanation, err = explainWithAI(errorInfo, level)
			if err == nil && store != nil {
				_ = store.Put(cache.Key(errorInfo, level, ""), explanation)
			}
		}

		if !explained && err != nil {
			fmt.Println()
			fmt.Println(
				"AIによる解説を取得できませんでした:",
				err,
			)
			return result.ExitCode
		}
	}

	policy := teaching.PolicyForLevel(
		level,
	)

	if policy.ShowSummary {
		fmt.Println()
		fmt.Println("Explanation")
		fmt.Println(explanation.Summary)
	}

	session := teaching.NewSession(
		explanation,
		policy,
	)

	terminal.RunHintSession(session)
	return result.ExitCode
}

func explainWithAI(
	errorInfo model.ErrorInfo,
	level teaching.Level,
) (model.Explanation, error) {
	provider, err := factory.NewProviderFromEnv(credential.NewKeyringStore())
	if err != nil {
		return model.Explanation{}, err
	}

	service := duckai.NewService(
		provider,
		redact.New(),
	)

	return service.Explain(
		context.Background(),
		errorInfo,
		level,
		"",
	)
}

func handleAuth(args []string) int {
	if len(args) == 1 && args[0] == "status" {
		handleGeminiAuthStatus(credential.NewKeyringStore())
		return 0
	}
	if len(args) != 2 || args[1] != "gemini" {
		fmt.Println("使い方: duck auth set|status|delete gemini")
		return 2
	}
	store := credential.NewKeyringStore()
	switch args[0] {
	case "set":
		fmt.Print("Gemini API key: ")
		secret, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			fmt.Println("APIキーを読み取れませんでした")
			return 1
		}
		if strings.TrimSpace(string(secret)) == "" {
			fmt.Println("APIキーが空です")
			return 2
		}
		if err := store.Set("DevDuck", "gemini", string(secret)); err != nil {
			fmt.Println("Credentialを保存できませんでした")
			return 1
		}
		fmt.Println("Credential saved securely.")
	case "status":
		handleGeminiAuthStatus(store)
	case "delete":
		err := store.Delete("DevDuck", "gemini")
		if err != nil && err != credential.ErrNotFound {
			fmt.Println("Credentialを削除できませんでした")
			return 1
		}
		fmt.Println("Credential deleted.")
	default:
		fmt.Println("使い方: duck auth set|status|delete gemini")
		return 2
	}
	return 0
}

func handleGeminiAuthStatus(store credential.Store) {
	_, err := store.Get("DevDuck", "gemini")
	if err == nil {
		fmt.Println("gemini: configured")
		return
	}
	fmt.Println("gemini: not configured")
}

func handleConfig(args []string) int {
	if len(args) == 0 {
		fmt.Println("使い方:")
		fmt.Println("duck config show")
		fmt.Println(
			"duck config set level <beginner|intermediate|advanced>",
		)
		return 2
	}

	switch args[0] {
	case "show":
		return handleConfigShow()

	case "set":
		return handleConfigSet(args[1:])

	default:
		fmt.Println(
			"不明なconfigコマンド:",
			args[0],
		)
		return 2
	}
}

func handleConfigShow() int {
	cfg, err := config.Load()

	if err != nil {
		fmt.Println(
			"設定の読み込みに失敗しました:",
			err,
		)
		return 1
	}

	path, err := config.Path()

	if err != nil {
		fmt.Println(
			"設定ファイルの場所を取得できません:",
			err,
		)
		return 1
	}

	fmt.Println("Level :", cfg.Level)
	fmt.Println("Config:", path)
	return 0
}

func handleConfigSet(args []string) int {
	if len(args) != 2 {
		fmt.Println(
			"使い方: duck config set level <beginner|intermediate|advanced>",
		)
		return 2
	}

	key := args[0]
	value := args[1]

	if key != "level" {
		fmt.Println(
			"未対応の設定項目:",
			key,
		)
		return 2
	}

	level, err := teaching.ParseLevel(
		value,
	)

	if err != nil {
		fmt.Println(
			"levelは beginner, intermediate, advanced のいずれかを指定してください。",
		)
		return 2
	}

	cfg, err := config.Load()

	if err != nil {
		fmt.Println(
			"設定の読み込みに失敗しました:",
			err,
		)
		return 1
	}

	cfg.Level = level

	if err := config.Save(cfg); err != nil {
		fmt.Println(
			"設定の保存に失敗しました:",
			err,
		)
		return 1
	}

	fmt.Println(
		"levelを",
		level,
		"に設定しました。",
	)
	return 0
}
