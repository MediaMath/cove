package gocmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type GoCmd struct {
	*exec.Cmd
}

func Go(sub string, args ...string) *GoCmd {
	arguments := append([]string{sub}, args...)
	return &GoCmd{exec.Command("go", arguments...)}
}

func (cmd *GoCmd) Receive(action func(io.Reader) error) error {
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

	return NewGoError(cmd.Wait(), NewStdError(stderr))
}

func (cmd *GoCmd) StdOutLines() ([]string, error) {
	results := make([]string, 0)
	err := cmd.Receive(func(stdout io.Reader) error {
		var scanerr error
		results, scanerr = scanReader(stdout)
		return scanerr
	})

	return results, err
}

func scanReader(reader io.Reader) ([]string, error) {
	results := make([]string, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	return results, scanner.Err()
}

type GoError struct {
	Exit   *exec.ExitError
	StdErr *StdError
}

func (s *GoError) Error() string {
	return fmt.Sprintf("%v\n%v", s.Exit, s.StdErr)
}

func NewGoError(err error, stdErr *StdError) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return &GoError{exitErr, stdErr}
	}

	return err
}

type StdError struct {
	Output string
}

func (s *StdError) Error() string {
	return s.Output
}

func NewStdError(stderr io.Reader) *StdError {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stderr)
	s := buf.String()

	if s != "" {
		return &StdError{s}
	} else {
		return nil
	}
}
