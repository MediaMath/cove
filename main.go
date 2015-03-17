package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	wd, _ := os.Getwd()

	fmt.Printf("go go go %v\n", wd)
	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if isPackage(path, info) {
			fmt.Printf("%v\n", info.Name())
		}

		return nil
	})
}

func isPackage(path string, info os.FileInfo) bool {
	return info.IsDir() && hasGoFiles(path)
}

func hasGoFiles(path string) bool {
	paths, readErr := ioutil.ReadDir(path)
	if readErr != nil {
		return false
	}

	for _, path := range paths {
		if isGoFile(path) {
			return true
		}
	}

	return false
}

func isGoFile(file os.FileInfo) bool {
	return filepath.Ext(file.Name()) == ".go"
}
