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
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
	"github.com/duck-ahiru-Z/DevDuck/internal/ui/terminal"
)

func main() {
	// DevDuck自身のオプションと、実行対象のコマンドを解析する。
	options, err := cli.Parse(os.Args[1:])

	if err != nil {
		fmt.Println("DevDuck:", err)
		fmt.Println()
		fmt.Println(
			"使い方: duck [--level beginner|intermediate|advanced] <command> [args...]",
		)
		return
	}

	command := options.Command
	args := options.Args

	// 対応しているAdapterを登録する。
	registry := adapter.NewRegistry(
		python.NewAdapter(),
	)

	// 実行対象に対応するAdapterを探す。
	selectedAdapter, adapterFound := registry.Find(
		command,
		args,
	)

	// ユーザーが指定したコマンドを実行する。
	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	// stderrは通常通り表示しながら、DevDuck側にも保存する。
	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(
		os.Stderr,
		&stderr,
	)

	err = cmd.Run()

	// 正常終了した場合、DevDuckは何も追加表示しない。
	if err == nil {
		return
	}

	fmt.Println()
	fmt.Println("DevDuck detected an error")

	// まだ対応していないコマンドだった場合。
	if !adapterFound {
		fmt.Println()
		fmt.Println("このコマンドにはまだ対応していません。")
		return
	}

	// stderrを解析する。
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

	// ユーザーのレベルに応じて教え方を変更する。
	policy := teaching.PolicyForLevel(
		options.Level,
	)

	if policy.ShowSummary {
		fmt.Println()
		fmt.Println("Explanation")
		fmt.Println(explanation.Summary)
	}

	// 段階的にヒントを表示する。
	session := teaching.NewSession(
		explanation,
		policy,
	)

	terminal.RunHintSession(session)
}
