package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/MediaMath/cove/list"
	"github.com/MediaMath/cove/test"
	"github.com/MediaMath/cove/tool"
	"github.com/pkg/browser"
)

func main() {
	outputDir := flag.String("o", "", "If provided output files will be written to this dir and not opened.")
	short := flag.Bool("short", false, "Run only short tests.")
	keepProfile := flag.Bool("keep", false, "Will not remove coverage profile files if set.")
	flag.Parse()

	openReport := false
	reportPath := *outputDir
	if reportPath == "" {
		openReport = true
		reportPath, _ = ioutil.TempDir("", "cvr")
	}

	paths := flag.Args()

	if len(paths) == 0 {
		paths = append(paths, ".")
	}

	packs, pathErr := list.Packages(paths...)
	logError(pathErr)

	profiles, coverErr := test.CoverageProfile(*short, reportPath, packs...)
	logError(coverErr)

	if len(profiles) < 1 {
		log.Printf("No coverage for %s, %s", paths, profiles)
	}

	for _, profile := range profiles {
		if !*keepProfile {
			defer os.RemoveAll(profile)
		}

		report, reportErr := tool.CoverageReport(profile, reportPath)
		logError(reportErr)

		if report != "" && openReport {
			browser.OpenFile(report)
		}

	}
}

func logError(err error) {
	if err != nil {
		log.Printf("%s\n", err)
	}
}
