package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/duck-ahiru-Z/DevDuck/internal/adapters/python"
	"github.com/duck-ahiru-Z/DevDuck/internal/cli"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
	"github.com/duck-ahiru-Z/DevDuck/internal/ui/terminal"
)

func main() {
	fmt.Println("DevDuck")

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

	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(
		os.Stderr,
		&stderr,
	)

	err = cmd.Run()

	if err != nil {
		fmt.Println("\n🦆 DevDuck detected an error")

		errorInfo, ok := python.Parse(stderr.String())

		if !ok {
			fmt.Println("エラーを解析できませんでした。")
			return
		}

		fmt.Println()
		fmt.Println("Source :", errorInfo.Source)
		fmt.Println("Type   :", errorInfo.Kind)
		fmt.Println("Message:", errorInfo.Message)
		fmt.Println("File   :", errorInfo.File)
		fmt.Println("Line   :", errorInfo.Line)

		explainer := python.NewExplainer()

		explanation, explained := explainer.Explain(errorInfo)

		if !explained {
			fmt.Println()
			fmt.Println(
				"このエラーはまだローカル解説に対応していません。",
			)
			return
		}

		policy := teaching.PolicyForLevel(options.Level)

		if policy.ShowSummary {
			fmt.Println()
			fmt.Println("🦆 Explanation")
			fmt.Println(explanation.Summary)
		}

		session := teaching.NewSession(
			explanation,
			policy,
		)

		terminal.RunHintSession(session)
	}
}
