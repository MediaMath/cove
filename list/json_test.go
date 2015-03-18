package list

import "fmt"

func ExampleJson() {
	var info struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	Json("os/exec", &info)

	fmt.Printf("%s %s %v", info.ImportPath, info.Name, info.Incomplete)

	//Output:
	//os/exec exec false
}

func ExampleJson_wildcard() {

	type JsonResponse struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	resp := make([]JsonResponse, 0)
	if err := Json("os/...", &resp); err != nil {
		fmt.Printf("'go list -json os/...' returns an invalid json in the multiple return case.")
	}

	//Output:
	//'go list -json os/...' returns an invalid json in the multiple return case.
}

func ExampleJson_unknown() {
	var info struct {
		ImportPath string
		Name       string
		Incomplete bool
	}

	if err := Json("moo/boo", &info); err != nil {
		fmt.Printf("'go list -json moo/boo' fails.")
	}

	//Output:
	//'go list -json moo/boo' fails.
}
