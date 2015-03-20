package cove

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func ExampleCoverageProfile() {
	testdir, temperr := ioutil.TempDir("", "coverage-profile")

	if temperr != nil {
		fmt.Printf("%v", temperr)
	}

	defer os.RemoveAll(testdir)

	profile, profileErr := CoverageProfile(false, "count", testdir, "text/scanner")
	if profileErr != nil {
		fmt.Printf("%v", profileErr)
	}

	fmt.Println(filepath.Base(profile))

	//Output:
	//text.scanner.out
}

func ExampleCoverageProfile_short() {
	testdir, _ := ioutil.TempDir("", "coverage-profile")
	defer os.RemoveAll(testdir)

	profile, _ := CoverageProfile(true, "set", testdir, "text/scanner")
	fmt.Println(filepath.Base(profile))

	//Output:
	//text.scanner.out
}
func ExamplePackages_single() {

	packs, _ := Packages("text/template")

	for _, pack := range packs {
		fmt.Println(pack)
	}

	//Output:
	//text/template
}

func ExamplePackages_wildcard() {
	packs, _ := Packages("text/...")

	for _, pack := range packs {
		fmt.Println(pack)
	}

	//Output:
	//text/scanner
	//text/tabwriter
	//text/template
	//text/template/parse
}

func ExamplePackages_multiples() {
	packs, _ := Packages("text/...", "os/exec")

	for _, pack := range packs {
		fmt.Println(pack)
	}

	//Output:
	//os/exec
	//text/scanner
	//text/tabwriter
	//text/template
	//text/template/parse
}

func ExamplePackages_relative() {
	//duplicates are filtered
	packs, _ := Packages(".")

	for _, pack := range packs {
		fmt.Println(pack)
	}

	//Output:
	//github.com/MediaMath/cove
}

func ExamplePackages_unknown() {
	packs, err := Packages("foo/manchu", "text/...", "moo/...")

	for _, pack := range packs {
		fmt.Println(pack)
	}

	if err != nil {
		fmt.Println("Had errors.")
	}

	//Output:
	//text/scanner
	//text/tabwriter
	//text/template
	//text/template/parse
	//Had errors.
}

func ExamplePackages_unknown2() {
	packs, err := Packages("foo/manchu")

	for _, pack := range packs {
		fmt.Println(pack)
	}

	if err != nil {
		fmt.Println("Had errors.")
	}

	//Output:
	//Had errors.
}

func ExamplePackageJSON() {
	var info struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	PackageJSON("os/exec", &info)

	fmt.Printf("%s %s %v", info.ImportPath, info.Name, info.Incomplete)

	//Output:
	//os/exec exec false
}

func ExamplePackageJSON_wildcard() {

	type JSONresponse struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	var resp []JSONresponse
	if err := PackageJSON("os/...", &resp); err != nil {
		fmt.Printf("'go list -json os/...' returns an invalid json in the multiple return case.")
	}

	//Output:
	//'go list -json os/...' returns an invalid json in the multiple return case.
}

func ExamplePackageJSON_unknown() {
	var info struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	if err := PackageJSON("moo/boo", &info); err != nil {
		fmt.Printf("'go list -json moo/boo' fails.")
	}

	//Output:
	//'go list -json moo/boo' fails.
}
