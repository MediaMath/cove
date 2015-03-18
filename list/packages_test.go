package list

import (
	"fmt"
	"sort"
)

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

	sort.Strings(packs)
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

	sort.Strings(packs)
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

	sort.Strings(packs)
	for _, pack := range packs {
		fmt.Println(pack)
	}

	//Output:
	//github.com/MediaMath/cove/list
}

func ExamplePackages_unknown() {
	packs, err := Packages("foo/manchu", "text/...", "moo/...")

	sort.Strings(packs)
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

	sort.Strings(packs)
	for _, pack := range packs {
		fmt.Println(pack)
	}

	if err != nil {
		fmt.Println("Had errors.")
	}

	//Output:
	//Had errors.
}
