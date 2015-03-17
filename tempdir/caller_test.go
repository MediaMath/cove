package tempdir

import (
	"regexp"
	"strings"
	"testing"
)

func TestGetTopCaller(t *testing.T) {
	file, line := getTopCaller()

	if !strings.Contains(file, "caller_test.go") {
		t.Errorf("Did not get caller correctly: %v %v", file, line)
	}
}

func TestGetCallerLabel(t *testing.T) {

	file, line := getTopCaller()
	label := getCallerLabel()

	match, _ := regexp.MatchString("caller_test.go-[0-9]+", label)
	if !match {
		t.Errorf("Bad label %v, %v %v", label, file, line)
	}

}
