package list

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func Packages(paths ...string) ([]string, error) {
	cmd, err := listResponse(execList(paths...))
	return cmd, err
}

func execList(args ...string) *exec.Cmd {
	arguments := append([]string{"list"}, args...)
	return exec.Command("go", arguments...)
}

func scanReader(reader io.Reader) ([]string, error) {
	results := make([]string, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	return results, scanner.Err()
}

func listResponse(cmd *exec.Cmd) ([]string, error) {
	results := make([]string, 0)

	stdout, stdOutErr := cmd.StdoutPipe()
	if stdOutErr != nil {
		return results, stdOutErr
	}

	stderr, stdErrErr := cmd.StderrPipe()
	if stdErrErr != nil {
		return results, stdErrErr
	}

	if err := cmd.Start(); err != nil {
		return results, err
	}

	results, scanOutErr := scanReader(stdout)
	if scanOutErr != nil {
		return results, scanOutErr
	}

	return results, NewExecError(cmd.Wait(), NewStdError(stderr))
}

type ExecError struct {
	Exit   *exec.ExitError
	StdErr *StdError
}

func (s *ExecError) Error() string {
	return fmt.Sprintf("%v\n%v", s.Exit, s.StdErr)
}

func NewExecError(err error, stdErr *StdError) error {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return &ExecError{exitErr, stdErr}
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

func jsonResponse(cmd *exec.Cmd, v interface{}) error {
	stdout, stdOutErr := cmd.StdoutPipe()
	if stdOutErr != nil {
		return stdOutErr
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	if err := json.NewDecoder(stdout).Decode(v); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func filterEmpty(sl []string) []string {
	filtered := make([]string, 0)
	for _, s := range sl {
		if strings.TrimSpace(s) != "" {
			filtered = append(filtered, s)
		}
	}

	return filtered
}
