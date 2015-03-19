package cove

import (
	"bufio"
	"fmt"
	"io"
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

func ExamplePrepare_receive() {
	Prepare("list", "os/exec", "text/...").Receive(func(stdout io.Reader) error {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		return nil
	})

	//Output:
	//os/exec
	//text/scanner
	//text/tabwriter
	//text/template
	//text/template/parse
}

func ExamplePrepare_stdout() {
	out, _ := Prepare("list", "os/exec", "text/...").StdOutLines()

	for _, k := range out {
		fmt.Println(k)
	}

	//Output:
	//os/exec
	//text/scanner
	//text/tabwriter
	//text/template
	//text/template/parse
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

	type JsonResponse struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	resp := make([]JsonResponse, 0)
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
