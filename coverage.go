package main

import (
	"fmt"
	"os"
	"os/exec"
)

func buildProfile(testLocation string, outFile string) (bool, error) {
	fmt.Printf("%v:%v\n", testLocation, outFile)
	cmd := exec.Command("go", "test", testLocation, "-cover", fmt.Sprintf("-coverprofile=%s", outFile))
	cmd.Run()

	info, exists := os.Stat(outFile)
	return exists == nil && info.Size() > 0, nil
}

func buildReport(coverProfile string, output string) error {
	return exec.Command("go", "tool", "cover", fmt.Sprintf("-html=%s", coverProfile), "-o", output).Run()
}
