package terminal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

func RunHintSession(session *teaching.Session) {
	reader := bufio.NewReader(os.Stdin)

	hintNumber := 1

	hint, ok := session.NextHint()

	if !ok {
		return
	}

	fmt.Println()
	fmt.Printf("Hint %d\n", hintNumber)
	fmt.Println(hint)

	for session.HasNextHint() {
		fmt.Println()
		fmt.Println("[h] 次のヒント")
		fmt.Println("[q] 終了")
		fmt.Print("> ")

		input, err := reader.ReadString('\n')

		if err != nil {
			return
		}

		input = strings.ToLower(strings.TrimSpace(input))

		switch input {
		case "h":
			hint, ok := session.NextHint()

			if !ok {
				return
			}

			hintNumber++

			fmt.Println()
			fmt.Printf("Hint %d\n", hintNumber)
			fmt.Println(hint)

		case "q":
			return

		default:
			fmt.Println("h または q を入力してください。")

		}
	}

	fmt.Println()
	fmt.Println("これ以上のローカルヒントはありません。")
}
