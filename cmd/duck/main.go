package main

import (
	"fmt"
	"os"
	"os/exec"
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
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Println("DevDuck: command failed")
	}
}
