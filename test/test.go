// Package test provides a simplified facade around 'go test'
package test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MediaMath/cove/gocmd"
)

// Coverage Profile creates a cover profile file for all of the provided packages.
// The files are created in outdir.  The parameter short sets whether to run
// all tests or only the short ones.
// If a profile is able to be created its file name is returned.
func CoverageProfile(short bool, outdir string, packs ...string) ([]string, error) {
	var written []string
	for _, pack := range packs {
		file, err := coverageProfile(short, outdir, pack)
		if err != nil {
			return written, fmt.Errorf("%v:%v", err, written)
		}

		if file != "" {
			written = append(written, file)
		}
	}

	return written, nil
}

func coverageProfile(short bool, outdir string, pack string) (string, error) {
	profile := getFileName(outdir, pack)

	cmd := gocmd.Prepare("test", pack, fmt.Sprintf("-coverprofile=%s", profile), getShort(short))
	if _, err := cmd.StdOutLines(); err != nil {
		return "", fmt.Errorf("%s:%v", pack, err)
	}

	if _, err := os.Stat(profile); err != nil {
		return "", nil
	}

	return profile, nil
}

func getFileName(outdir string, pack string) string {
	profile := strings.Replace(pack, "/", ".", -1)
	fullPath := filepath.Join(outdir, profile)
	return fmt.Sprintf("%s.out", fullPath)
}

func getShort(short bool) string {
	if short {
		return "-short"
	} else {
		return ""
	}
}
