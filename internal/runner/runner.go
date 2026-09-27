package runner

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type Result struct {
	ExitCode int
	Stderr   string
	Err      error
}

func Run(command string, args []string, stdin io.Reader, stdout, stderr io.Writer) Result {
	cmd := exec.Command(command, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	var captured bytes.Buffer
	cmd.Stderr = io.MultiWriter(stderr, &captured)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return Result{ExitCode: exitErr.ProcessState.ExitCode(), Stderr: captured.String(), Err: err}
		}
		return Result{ExitCode: 1, Stderr: captured.String(), Err: err}
	}
	return Result{ExitCode: 0, Stderr: captured.String()}
}

func Streams() (io.Reader, io.Writer, io.Writer) { return os.Stdin, os.Stdout, os.Stderr }
