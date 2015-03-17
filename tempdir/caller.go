package tempdir

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func getTopCaller() (string, int) {

	var i int = 0

	for {
		_, _, _, ok := runtime.Caller(i)

		if !ok {
			break
		}

		i++
	}

	//take out the golang system calls in the stack
	//if you want to see them, remove this line
	i -= 3

	_, file, line, _ := runtime.Caller(i)

	return file, line
}

func getCallerLabel() string {
	file, line := getTopCaller()
	return fmt.Sprintf("%s-%v", filepath.Base(file), line)
}
