package cove

import (
	"bufio"
	"io"
	"os/exec"
	"strings"
)

func run(cmd *exec.Cmd) error {
	return pipeWith(cmd, func(stdout io.Reader) error { return nil })
}

func pipeWith(cmd *exec.Cmd, action func(io.Reader) error) error {
	stdout, stdOutErr := cmd.StdoutPipe()
	if stdOutErr != nil {
		return stdOutErr
	}

	stderr, stdErrErr := cmd.StderrPipe()
	if stdErrErr != nil {
		return stdErrErr
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := action(stdout); err != nil {
		return err
	}

	stderrLines, stderrLinesErr := scanReader(stderr)
	if stderrLinesErr != nil {
		return stderrLinesErr
	}

	err := newGoError(cmd.Wait(), stderrLines)
	return err
}

func output(cmd *exec.Cmd) ([]string, error) {
	var results []string
	err := pipeWith(cmd, func(stdout io.Reader) error {
		var scanerr error
		results, scanerr = scanReader(stdout)
		return scanerr
	})

	return results, err
}

func scanReader(reader io.Reader) ([]string, error) {
	var results []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	return results, scanner.Err()
}

//GoError is returned when a go command fails during its execution.
type GoError struct {
	Exit   *exec.ExitError
	StdErr []string
}

func (s *GoError) Error() string {
	stderr := strings.Join(s.StdErr, "\n")
	if stderr != "" {
		return stderr
	}

	return s.Exit.Error()
}

func newGoError(err error, stderr []string) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return &GoError{exitErr, stderr}
	}

	return err
}
