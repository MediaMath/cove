package gocmd

import (
	"bufio"
	"fmt"
	"io"
)

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
