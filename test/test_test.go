package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ExampleCoverageProfile() {
	testdir, _ := ioutil.TempDir("", "coverage-profile")
	defer os.RemoveAll(testdir)

	written, _ := CoverageProfile(false, testdir, "text/scanner", "text/tabwriter")

	for _, profile := range written {
		fmt.Println(filepath.Base(profile))
	}

	//Output:
	//text.scanner.out
	//text.tabwriter.out
}

func ExampleCoverageProfile_short() {
	testdir, _ := ioutil.TempDir("", "coverage-profile")
	defer os.RemoveAll(testdir)

	written, _ := CoverageProfile(true, testdir, "text/scanner", "text/tabwriter")

	for _, profile := range written {
		fmt.Println(filepath.Base(profile))
	}

	//Output:
	//text.scanner.out
	//text.tabwriter.out
}
