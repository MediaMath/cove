package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func TestCmdError(t *testing.T) {
	cmd, dir := getFailingCmd()
	err := Run(cmd)

	if err == nil {
		t.Errorf("Should have error")
	}

	if goerr, ok := err.(*CmdError); ok {
		expected := fmt.Sprintf("ls: %v: No such file or directory", dir)
		if goerr.Error() != expected {
			t.Errorf("%v\n%v", goerr.Error(), expected)
		}
	} else {
		t.Errorf("Should be CmdError: %v", err)
	}
}

func TestCmdErrorDoesntTrapNils(t *testing.T) {
	err := newCmdError(nil, []string{"foo"})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestCmdErrorDoesntTrapNonExitErrors(t *testing.T) {
	in := fmt.Errorf("foo")
	err := newCmdError(in, []string{"foo"})
	if err != in {
		t.Errorf("%v", err)
	}
}

//to get an exit err I need to run a process
//and get that error
func getRealExitError() error {
	cmd, _ := getFailingCmd()
	return cmd.Run()
}

func getFailingCmd() (*exec.Cmd, string) {
	tempDir, _ := ioutil.TempDir("", "failingcmd")
	os.RemoveAll(tempDir)

	//ls on a non-existent directory should fail
	//not portable to systems without ls
	return exec.Command("ls", tempDir), tempDir
}

func TestCmdErrorTrapsExitErrors(t *testing.T) {
	in := getRealExitError()
	err := newCmdError(in, []string{"foo"})
	if err == in {
		t.Errorf("%v", err)
	}
}

func TestCmdErrorUsesStdErr(t *testing.T) {
	in := getRealExitError()
	err := newCmdError(in, []string{"foo", "bar"})
	if err == in {
		t.Errorf("%v", err)
	}

	if err.Error() != "foo\nbar" {
		t.Errorf("%v", err)
	}
}

func TestCmdErrorUsesInErrorInFaceOfNilStdErr(t *testing.T) {
	in := getRealExitError()
	err := newCmdError(in, nil)

	if err.Error() != in.Error() {
		t.Errorf("%v", err)
	}
}
