package cli

import (
	"fmt"

	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

type Options struct {
	Level   teaching.Level
	Command string
	Args    []string
}

func Parse(args []string) (Options, error) {
	options := Options{
		Level: teaching.Intermediate,
	}

	index := 0

	for index < len(args) {
		arg := args[index]

		if arg == "--level" {
			if index+1 >= len(args) {
				return Options{}, fmt.Errorf(
					"--level の後にレベルを指定してください",
				)
			}

			level, err := teaching.ParseLevel(args[index+1])

			if err != nil {
				return Options{}, err
			}

			options.Level = level
			index += 2
			continue
		}

		// 最初のオプションではない値を
		// 実行コマンドとして扱う
		options.Command = arg
		options.Args = args[index+1:]

		return options, nil
	}

	return Options{}, fmt.Errorf("実行するコマンドが指定されていません")
}
