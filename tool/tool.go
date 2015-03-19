// Package tool provides a thin wrapper around 'go tool'
package tool

import (
	"fmt"
	"path/filepath"

	"github.com/MediaMath/cove/gocmd"
)

// CoverageReport turns the profile into a report using 'go tool cover'
func CoverageReport(profile string, outdir string) (string, error) {
	report := getFileName(profile, outdir)
	if _, err := gocmd.Prepare("tool", "cover", fmt.Sprintf("-html=%s", profile), "-o", report).StdOutLines(); err != nil {
		return "", err
	}

	return report, nil
}

func getFileName(profile string, outdir string) string {
	report := filepath.Base(profile)
	extension := filepath.Ext(report)
	name := report[0 : len(report)-len(extension)]
	fullPath := filepath.Join(outdir, name)
	return fmt.Sprintf("%s.html", fullPath)
}
