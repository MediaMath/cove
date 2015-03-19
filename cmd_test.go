package cove

import (
	"fmt"
	"testing"
)

func TestGoError(t *testing.T) {
	err := run(GoCmd(""))

	if err == nil {
		t.Errorf("Should have error")
	}

	if goerr, ok := err.(*GoError); ok {
		expected := "go: unknown subcommand \"\"\nRun 'go help' for usage."
		if goerr.Error() != expected {
			t.Errorf("%v\n%v", goerr.Error(), expected)
		}
	} else {
		t.Errorf("Should be GoError: %v", err)
	}
}

func TestGoErrorDoesntTrapNils(t *testing.T) {
	err := newGoError(nil, []string{"foo"})
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestGoErrorDoesntTrapNonExitErrors(t *testing.T) {
	in := fmt.Errorf("foo")
	err := newGoError(in, []string{"foo"})
	if err != in {
		t.Errorf("%v", err)
	}
}

//to get an exit err I need to run a process
//and get that error
func getRealExitError() error {
	return GoCmd("").Run()
}

func TestGoErrorTrapsExitErrors(t *testing.T) {
	in := getRealExitError()
	err := newGoError(in, []string{"foo"})
	if err == in {
		t.Errorf("%v", err)
	}
}

func TestGoErrorUsesStdErr(t *testing.T) {
	in := getRealExitError()
	err := newGoError(in, []string{"foo", "bar"})
	if err == in {
		t.Errorf("%v", err)
	}

	if err.Error() != "foo\nbar" {
		t.Errorf("%v", err)
	}
}

func TestGoErrorUsesInErrorInFaceOfNilStdErr(t *testing.T) {
	in := getRealExitError()
	err := newGoError(in, nil)

	if err.Error() != in.Error() {
		t.Errorf("%v", err)
	}
}
