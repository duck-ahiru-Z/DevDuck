package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/duck-ahiru-Z/DevDuck/internal/adapter"
	"github.com/duck-ahiru-Z/DevDuck/internal/adapters/python"
	duckai "github.com/duck-ahiru-Z/DevDuck/internal/ai"
	"github.com/duck-ahiru-Z/DevDuck/internal/ai/providers/gemini"
	"github.com/duck-ahiru-Z/DevDuck/internal/cache"
	"github.com/duck-ahiru-Z/DevDuck/internal/cli"
	"github.com/duck-ahiru-Z/DevDuck/internal/config"
	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/redact"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
	"github.com/duck-ahiru-Z/DevDuck/internal/ui/terminal"
)

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "config" {
		handleConfig(os.Args[2:])
		return
	}

	cfg, err := config.Load()

	if err != nil {
		fmt.Println("設定ファイルの読み込みに失敗しました:", err)
		return
	}

	options, err := cli.Parse(os.Args[1:])

	if err != nil {
		fmt.Println("DevDuck:", err)
		fmt.Println()
		fmt.Println(
			"使い方: duck [--level beginner|intermediate|advanced] <command> [args...]",
		)
		return
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

	cmd := exec.Command(
		command,
		args...,
	)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(
		os.Stderr,
		&stderr,
	)

	err = cmd.Run()

	if err == nil {
		return
	}

	fmt.Println()
	fmt.Println("DevDuck detected an error")

	if !adapterFound {
		fmt.Println()
		fmt.Println("このコマンドにはまだ対応していません。")
		return
	}

	errorInfo, ok := selectedAdapter.Parse(
		stderr.String(),
	)

	if !ok {
		fmt.Println()
		fmt.Println("エラーを解析できませんでした。")
		return
	}

	fmt.Println()
	fmt.Println("Source :", errorInfo.Source)
	fmt.Println("Type   :", errorInfo.Kind)
	fmt.Println("Message:", errorInfo.Message)
	fmt.Println("File   :", errorInfo.File)
	fmt.Println("Line   :", errorInfo.Line)

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
			return
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
}

func explainWithAI(
	errorInfo model.ErrorInfo,
	level teaching.Level,
) (model.Explanation, error) {
	apiKey := os.Getenv(
		"GEMINI_API_KEY",
	)

	if apiKey == "" {
		return model.Explanation{},
			fmt.Errorf(
				"GEMINI_API_KEY が設定されていません",
			)
	}

	modelName := os.Getenv(
		"DEVDUCK_GEMINI_MODEL",
	)

	if modelName == "" {
		modelName = "gemini-3.5-flash-lite"
	}

	provider := gemini.New(
		apiKey,
		modelName,
	)

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

func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Println("使い方:")
		fmt.Println("duck config show")
		fmt.Println(
			"duck config set level <beginner|intermediate|advanced>",
		)
		return
	}

	switch args[0] {
	case "show":
		handleConfigShow()

	case "set":
		handleConfigSet(args[1:])

	default:
		fmt.Println(
			"不明なconfigコマンド:",
			args[0],
		)
	}
}

func handleConfigShow() {
	cfg, err := config.Load()

	if err != nil {
		fmt.Println(
			"設定の読み込みに失敗しました:",
			err,
		)
		return
	}

	path, err := config.Path()

	if err != nil {
		fmt.Println(
			"設定ファイルの場所を取得できません:",
			err,
		)
		return
	}

	fmt.Println("Level :", cfg.Level)
	fmt.Println("Config:", path)
}

func handleConfigSet(args []string) {
	if len(args) != 2 {
		fmt.Println(
			"使い方: duck config set level <beginner|intermediate|advanced>",
		)
		return
	}

	key := args[0]
	value := args[1]

	if key != "level" {
		fmt.Println(
			"未対応の設定項目:",
			key,
		)
		return
	}

	level, err := teaching.ParseLevel(
		value,
	)

	if err != nil {
		fmt.Println(
			"levelは beginner, intermediate, advanced のいずれかを指定してください。",
		)
		return
	}

	cfg, err := config.Load()

	if err != nil {
		fmt.Println(
			"設定の読み込みに失敗しました:",
			err,
		)
		return
	}

	cfg.Level = level

	if err := config.Save(cfg); err != nil {
		fmt.Println(
			"設定の保存に失敗しました:",
			err,
		)
		return
	}

	fmt.Println(
		"levelを",
		level,
		"に設定しました。",
	)
}
