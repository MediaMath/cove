package cove

import (
	"reflect"
	"testing"
)

//These tests verify the behavior of go list
func TestSinglePathResult(t *testing.T) {
	cmd, err := Prepare("list", "os/exec").StdOutLines()
	if !eq(cmd, sl("os/exec")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestMultiPathResult(t *testing.T) {
	cmd, err := Prepare("list", "os/...").StdOutLines()
	if !eq(cmd, sl("os", "os/exec", "os/signal", "os/user")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResult(t *testing.T) {
	cmd, err := Prepare("list", "os/foodog").StdOutLines()
	if len(cmd) > 0 || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingToo(t *testing.T) {
	cmd, err := Prepare("list", "os/foodog", "os/signal").StdOutLines()
	if !eq(cmd, sl("os/signal")) || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithErrorsOff(t *testing.T) {
	cmd, err := Prepare("list", "-e", "os/foodog").StdOutLines()
	if !eq(cmd, sl("os/foodog")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingTooWithErrorsOff(t *testing.T) {
	cmd, err := Prepare("list", "-e", "os/foodog", "os/signal").StdOutLines()
	if !eq(cmd, sl("os/foodog", "os/signal")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func eq(slice1 []string, slice2 []string) bool {
	return reflect.DeepEqual(slice1, slice2)
}

func sl(items ...string) []string {
	return items
}
