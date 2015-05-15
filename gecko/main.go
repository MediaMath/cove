package main

import (
	"fmt"
	"os"
)

func main() {
	out := "darnit"
	if len(os.Args) == 2 {
		out = os.Args[1]
	}

	fmt.Printf("gosh %v\n", out)
}
