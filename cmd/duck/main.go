package main

import (
	"bytes"
	"fmt"
	"io"
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

	var stderr bytes.Buffer

	cmd.Stderr = io.MultiWriter(
		os.Stderr,
		&stderr,
	)

	err := cmd.Run()

	if err != nil {
		fmt.Println("\n DevDuck detected an error")
		fmt.Println("----- captured stderr -----")
		fmt.Println(stderr.String())
	}
}
