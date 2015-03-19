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

	err := cmd.Wait()
	gerr := newGoError(err, newStdError(stderrLines))
	return gerr
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
	StdErr *StdError
}

func (s *GoError) Error() string {
	if s.StdErr != nil {
		return s.StdErr.Error()
	}

	return s.Exit.Error()
}

func newGoError(err error, stdErr *StdError) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return &GoError{exitErr, stdErr}
	}

	return err
}

//StdError is returned anytime a go command prints to stderr during execution.
type StdError struct {
	Output string
}

func (s *StdError) Error() string {
	return s.Output
}

func newStdError(lines []string) *StdError {
	s := strings.Join(lines, "\n")

	if s != "" {
		return &StdError{s}
	}

	return nil
}
