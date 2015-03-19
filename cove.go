// Package provides library wrappings around the go toolchain
package cove

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Packages gets all packages that match any of the paths.
// The package list will only contain 1 entry per package in sorted order.
// Invalid paths will generate an error, but will not stop the evaluation of the other paths.
func Packages(paths ...string) ([]string, error) {
	packs, err := Prepare("list", paths...).StdOutLines()
	sort.Strings(packs)
	return packs, err
}

// PackageJSON takes a SINGLE fully qualified package import path and decodes the 'go list -json' response.
// See $GOROOT/src/cmd/go/list.go for documentation on the json output.
func PackageJSON(pack string, v interface{}) error {
	return Prepare("list", "-json", pack).Receive(func(stdout io.Reader) error {
		return json.NewDecoder(stdout).Decode(v)
	})
}

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

// CoverageReport turns the profile into a report using 'go tool cover'
func CoverageReport(profile string, outdir string) (string, error) {
	report := getReportFileName(profile, outdir)
	if _, err := Prepare("tool", "cover", fmt.Sprintf("-html=%s", profile), "-o", report).StdOutLines(); err != nil {
		return "", err
	}

	return report, nil
}

func getReportFileName(profile string, outdir string) string {
	report := filepath.Base(profile)
	extension := filepath.Ext(report)
	name := report[0 : len(report)-len(extension)]
	fullPath := filepath.Join(outdir, name)
	return fmt.Sprintf("%s.html", fullPath)
}
func coverageProfile(short bool, outdir string, pack string) (string, error) {
	profile := getProfileFileName(outdir, pack)

	cmd := Prepare("test", pack, fmt.Sprintf("-coverprofile=%s", profile), getShort(short))
	if _, err := cmd.StdOutLines(); err != nil {
		return "", fmt.Errorf("%s:%v", pack, err)
	}

	if _, err := os.Stat(profile); err != nil {
		return "", nil
	}

	return profile, nil
}

func getProfileFileName(outdir string, pack string) string {
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
