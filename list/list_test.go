package list

import (
	"reflect"
	"testing"
)

func TestSinglePathResult(t *testing.T) {
	cmd, err := listOutput(execList("os/exec"))
	if !eq(cmd, sl("os/exec")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestMultiPathResult(t *testing.T) {
	cmd, err := listOutput(execList("os/..."))
	if !eq(cmd, sl("os", "os/exec", "os/signal", "os/user")) {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResult(t *testing.T) {
	cmd, err := listOutput(execList("os/foodog"))
	if len(cmd) > 0 || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingToo(t *testing.T) {
	cmd, err := listOutput(execList("os/foodog", "os/signal"))
	if !eq(cmd, sl("os/signal")) || err == nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithErrorsOff(t *testing.T) {
	cmd, err := listOutput(execList("-e", "os/foodog"))
	if !eq(cmd, sl("os/foodog")) || err != nil {
		t.Errorf("|%v|%v|", cmd, err)
	}
}

func TestFailPathResultWithPassingTooWithErrorsOff(t *testing.T) {
	cmd, err := listOutput(execList("-e", "os/foodog", "os/signal"))
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

func TestQueryParse(t *testing.T) {
	path, dir, err := queryParse("'text/template|/usr/local/Cellar/go/1.4.2/libexec/src/text/template|<nil>'")

	if path != "text/template" {
		t.Errorf("%v|%v|%v", path, dir, err)
	}

	if dir != "/usr/local/Cellar/go/1.4.2/libexec/src/text/template" {
		t.Errorf("%v|%v|%v", path, dir, err)
	}

	if err != nil {
		t.Errorf("%v|%v|%v", path, dir, err)
	}
}

func TestQueryParseFailure(t *testing.T) {
	path, dir, err := queryParse("'text/goog||package text/goog: cannot find package \"text/goog\" in any of:       /usr/local/Cellar/go/1.4.2/libexec/src/text/goog (from $GOROOT)         //go/src/t    ext/goog (from $GOPATH)'")

	if path != "text/goog" {
		t.Errorf("%v|%v|%v", path, dir, err)
	}

	if dir != "" {
		t.Errorf("%v|%v|%v", path, dir, err)
	}

	if err == nil {
		t.Errorf("%v|%v|%v", path, dir, err)
	}
}

func TestListQueryTemplateSingle(t *testing.T) {
	expected := sl("'text/template|/usr/local/Cellar/go/1.4.2/libexec/src/text/template|<nil>'")
	mp, errs := listQueryTemplate("text/template")
	if !eq(expected, mp) {
		t.Errorf("%v|%v|%v", expected, mp, errs)
	}
}

func TestQuerySingle(t *testing.T) {
	md, errs := Query("text/template")
	if _, contains := md["text/template"]; !contains || len(errs) > 0 {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestQueryMultiple(t *testing.T) {
	md, errs := Query("text/...")
	if len(errs) > 0 {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/scanner"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/tabwriter"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/template"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/template/parse"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestQueryMultipleWithDupes(t *testing.T) {
	md, errs := Query("text/...", "text/template")
	if len(errs) > 0 {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/template"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/template/parse"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestQueryFailure(t *testing.T) {
	md, errs := Query("text/foomanchu")

	if _, contains := md["text/foomanchu"]; contains {
		t.Errorf("%v|%v", md, errs)
	}

	if len(errs) < 1 {
		t.Errorf("%v|%v", md, errs)
	}
}

func TestQueryMixedSuccessAndMultipleFailures(t *testing.T) {
	md, errs := Query("text/foomanchu", "text/...", "moo/boo")

	if len(errs) < 2 {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/foomanchu"]; contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["moo/boo"]; contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/template"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}

	if _, contains := md["text/template/parse"]; !contains {
		t.Errorf("%v|%v", md, errs)
	}
}
