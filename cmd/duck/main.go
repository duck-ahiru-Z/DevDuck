package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/duck-ahiru-Z/DevDuck/internal/adapter"
	"github.com/duck-ahiru-Z/DevDuck/internal/adapters/python"
	"github.com/duck-ahiru-Z/DevDuck/internal/cli"
	"github.com/duck-ahiru-Z/DevDuck/internal/config"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
	"github.com/duck-ahiru-Z/DevDuck/internal/ui/terminal"
)

func main() {
	// config コマンドを処理する。
	if len(os.Args) >= 2 &&
		os.Args[1] == "config" {

		handleConfig(os.Args[2:])
		return
	}

	// 保存済み設定を読み込む。
	cfg, err := config.Load()

	if err != nil {
		fmt.Println(
			"設定ファイルの読み込みに失敗しました:",
			err,
		)
		return
	}

	// CLIオプションを解析する。
	options, err := cli.Parse(os.Args[1:])

	if err != nil {
		fmt.Println("DevDuck:", err)
		fmt.Println()
		fmt.Println(
			"使い方: duck [--level beginner|intermediate|advanced] <command> [args...]",
		)
		return
	}

	// --level が指定されていれば、
	// 保存済み設定より優先する。
	level := cfg.Level

	if options.LevelSet {
		level = options.Level
	}

	command := options.Command
	args := options.Args

	// 対応Adapterを登録する。
	registry := adapter.NewRegistry(
		python.NewAdapter(),
	)

	selectedAdapter, adapterFound :=
		registry.Find(command, args)

	// 指定されたコマンドを実行する。
	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(
		os.Stderr,
		&stderr,
	)

	err = cmd.Run()

	// 成功した場合は追加処理をしない。
	if err == nil {
		return
	}

	fmt.Println()
	fmt.Println("DevDuck detected an error")

	if !adapterFound {
		fmt.Println()
		fmt.Println(
			"このコマンドにはまだ対応していません。",
		)
		return
	}

	// stderrを解析する。
	errorInfo, ok :=
		selectedAdapter.Parse(
			stderr.String(),
		)

	if !ok {
		fmt.Println()
		fmt.Println(
			"エラーを解析できませんでした。",
		)
		return
	}

	fmt.Println()
	fmt.Println("Source :", errorInfo.Source)
	fmt.Println("Type   :", errorInfo.Kind)
	fmt.Println("Message:", errorInfo.Message)
	fmt.Println("File   :", errorInfo.File)
	fmt.Println("Line   :", errorInfo.Line)

	// ローカル解説を取得する。
	explanation, explained :=
		selectedAdapter.Explain(errorInfo)

	if !explained {
		fmt.Println()
		fmt.Println(
			"このエラーはまだローカル解説に対応していません。",
		)
		return
	}

	policy := teaching.PolicyForLevel(level)

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

func handleConfig(args []string) {
	if len(args) == 0 {
		fmt.Println(
			"使い方: duck config show",
		)
		fmt.Println(
			"       duck config set level <level>",
		)
		return
	}

	switch args[0] {
	case "show":
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

	case "set":
		handleConfigSet(args[1:])

	default:
		fmt.Println(
			"不明なconfigコマンド:",
			args[0],
		)
	}
}

func handleConfigSet(args []string) {
	if len(args) != 2 {
		fmt.Println(
			"使い方: duck config set level <level>",
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

	level, err := teaching.ParseLevel(value)

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
