// Package tool provides a thin wrapper around 'go tool'
package tool

import (
	"fmt"

	"github.com/MediaMath/cove/gocmd"
)

// CoverageReport turns the profile into a report using 'go tool cover'
func CoverageReport(profile string, report string) error {
	_, err := gocmd.Prepare("tool", "cover", fmt.Sprintf("-html=%s", profile), "-o", report).StdOutLines()
	return err
}
