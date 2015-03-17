package tempdir

import (
	"os"
	"testing"
)

func TestTempDirCreatesDirectory(t *testing.T) {
	In(func(dir string) {
		stat, err := os.Stat(dir)
		if err != nil {
			t.Errorf("Error on stats of temp dir: %v", err)
		}

		if !stat.IsDir() {
			t.Errorf("Tempdir is not a directroy: %v", stat)
		}
	})
}

func TestTempDirDeletesItself(t *testing.T) {
	var created string
	In(func(dir string) {
		created = dir
	})

	_, err := os.Stat(created)
	if err == nil {
		t.Errorf("File exists.")
	}
}
