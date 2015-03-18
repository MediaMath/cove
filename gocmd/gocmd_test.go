package gocmd

import (
	"bufio"
	"fmt"
	"io"
)

func ExampleGoCmdReceive() {
	Go("list", "os/exec", "text/...").Receive(func(stdout io.Reader) error {
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

func ExampleGoCmdStdOutLines() {
	out, _ := Go("list", "os/exec", "text/...").StdOutLines()

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
