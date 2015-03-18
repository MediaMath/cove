package list

import (
	"reflect"
	"testing"
)

type JsonResponse struct {
	ImportPath string
	Name       string
	Incomplete bool
}

func TestSinglePathResultJson(t *testing.T) {
	var info JsonResponse

	if err := jsonResponse(execList("-json", "os/exec"), &info); err != nil {
		t.Fatal(err)
	}

	if info.ImportPath != "os/exec" || info.Name != "exec" || info.Incomplete != false {
		t.Errorf("%v", info)
	}
}

func TestMultiplePathResultJson(t *testing.T) {
	//go list returns invalid json in the multiple return case.  Thanks Obama!
	resp := make([]JsonResponse, 0)
	if err := jsonResponse(execList("-json", "os/..."), &resp); err == nil {
		t.Errorf("Should have failed.")
	}
}

func TestSingleUnknownResultJson(t *testing.T) {
	var info JsonResponse
	if err := jsonResponse(execList("-e", "-json", "os/moo"), &info); err != nil {
		t.Fatal(err)
	}

	if info.Incomplete != true {
		t.Errorf("%v", info)
	}
}

func TestSinglePathResult(t *testing.T) {
	cmd, err := listResponse(execList("os/exec"))
	if !eq(cmd, sl("os/exec")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestMultiPathResult(t *testing.T) {
	cmd, err := listResponse(execList("os/..."))
	if !eq(cmd, sl("os", "os/exec", "os/signal", "os/user")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResult(t *testing.T) {
	cmd, err := listResponse(execList("os/foodog"))
	if len(cmd) > 0 || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingToo(t *testing.T) {
	cmd, err := listResponse(execList("os/foodog", "os/signal"))
	if !eq(cmd, sl("os/signal")) || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithErrorsOff(t *testing.T) {
	cmd, err := listResponse(execList("-e", "os/foodog"))
	if !eq(cmd, sl("os/foodog")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingTooWithErrorsOff(t *testing.T) {
	cmd, err := listResponse(execList("-e", "os/foodog", "os/signal"))
	if !eq(cmd, sl("os/foodog", "os/signal")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestPackagesSingle(t *testing.T) {
	md, errs := Packages("text/template")

	if len(md) != 1 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/template") != 1 {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestPackagesMultiple(t *testing.T) {
	md, errs := Packages("text/...")

	if len(md) != 4 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/scanner") != 1 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/tabwriter") != 1 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/template") != 1 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/template/parse") != 1 {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestPackagesMultipleWithDupes(t *testing.T) {
	md, errs := Packages("text/...", "text/template")

	if len(md) != 4 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/template") != 1 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/template/parse") != 1 {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestPackagesFailure(t *testing.T) {
	md, err := Packages("text/foomanchu")

	if len(md) != 0 {
		t.Errorf("%v|%v", md, err)
	}

	if err == nil {
		t.Errorf("%v|%v", md, err)
	}
}

func TestPackagesMixedSuccessAndMultipleFailures(t *testing.T) {
	md, errs := Packages("text/foomanchu", "text/...", "moo/boo")

	if errs == nil {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "text/foomanchu") != 0 {
		t.Errorf("%v|%v", md, errs)
	}

	if count(md, "moo/boo") != 0 {
		t.Errorf("%v|%v", md, errs)
	}

	if len(md) != 4 {
		t.Errorf("%v|%v", md, errs)
	}
}

func count(slice []string, item string) int {
	c := 0
	for _, s := range slice {
		if s == item {
			c = c + 1
		}
	}
	return c
}

func eq(slice1 []string, slice2 []string) bool {
	return reflect.DeepEqual(slice1, slice2)
}

func sl(items ...string) []string {
	return items
}
