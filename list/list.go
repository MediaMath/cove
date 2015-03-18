package list

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

func Query(paths ...string) (map[string]string, []error) {
	cmd, err := listQueryTemplate(paths...)
	if err != nil {
		return map[string]string{}, []error{err}
	}

	errors := make([]error, 0)
	results := make(map[string]string)
	for _, res := range cmd {
		path, dir, parseErr := queryParse(res)
		if parseErr != nil {
			errors = append(errors, parseErr)
		} else {
			results[path] = dir
		}
	}

	return results, errors
}

func listQueryTemplate(paths ...string) ([]string, error) {
	args := append([]string{"-e", "-f", "'{{ .ImportPath }}|{{ .Dir }}|{{ .Error }}'"}, paths...)
	return listOutput(execList(args...))
}

func queryParse(result string) (string, string, error) {
	result = strings.TrimPrefix(result, "'")
	result = strings.TrimSuffix(result, "'")
	split := strings.Split(result, "|")
	if len(split) != 3 {
		return "", "", fmt.Errorf("Unable to parse result: %v", result)
	}

	errString := split[2]
	var err error
	if errString != "<nil>" {
		err = fmt.Errorf(errString)
	}

	return split[0], split[1], err
}

func execList(args ...string) *exec.Cmd {
	arguments := append([]string{"list"}, args...)
	return exec.Command("go", arguments...)
}

func listOutput(cmd *exec.Cmd) ([]string, error) {
	out, err := cmd.Output()
	return filterEmpty(strings.Split(strings.TrimSpace(fmt.Sprintf("%s", out)), "\n")), err
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
