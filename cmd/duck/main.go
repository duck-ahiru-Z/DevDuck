package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/duck-ahiru-Z/DevDuck/internal/adapters/python"
)

func main() {
	fmt.Println("DevDuck")

	if len(os.Args) < 2 {
		fmt.Println("使い方: duck <command> [args...]")
		return
	}

	command := os.Args[1]
	args := os.Args[2:]

	cmd := exec.Command(command, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(
		os.Stderr,
		&stderr,
	)

	err := cmd.Run()

	if err != nil {
		fmt.Println("\nDevDuck detected an error")

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
			fmt.Println("このエラーはまだローカル解説に対応していません。")
			return
		}

		fmt.Println()
		fmt.Println("Explanation")
		fmt.Println(explanation.Summary)

		if len(explanation.Hints) > 0 {
			fmt.Println()
			fmt.Println("Hint:")
			fmt.Println(explanation.Hints[0])
		}
	}
}
