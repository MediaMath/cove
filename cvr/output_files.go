package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type OutputFiles interface {
	profile(string) string
	html(string) string
}

type existingDir string

func TempOutputFiles() (OutputFiles, error) {
	tempDir, err := ioutil.TempDir("", "cvr")
	if err != nil {
		return nil, err
	}

	return existingDir(tempDir), nil
}

func ExistingDirOutputFiles(path string) (OutputFiles, error) {
	if err := os.MkdirAll(path, 0644); err != nil {
		return nil, err
	}

	return existingDir(path), nil
}

func (dir existingDir) html(flatName string) string {
	return dir.file(addExtension(flatName, "html"))
}

func (dir existingDir) profile(flatName string) string {
	return dir.file(addExtension(flatName, "out"))
}

func (dir existingDir) file(path string) string {
	return filepath.Join(string(dir), path)
}

func addExtension(location string, extension string) string {
	return fmt.Sprintf("%s.%s", location, extension)
}