package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pkg/browser"
)

func main() {
	outputDir := flag.String("o", "", "If provided output files will be written to this dir and not opened.")
	short := flag.Bool("short", false, "Run only short tests.")
	keepProfile := flag.Bool("keep", false, "Will not remove coverage profile files if set.")
	flag.Parse()

	config, confErr := Config(*outputDir, *short, *keepProfile)
	if confErr != nil {
		log.Fatalf("Configuration Error: %v", confErr)
	}

	paths, pathErr := toPaths(flag.Args())
	if pathErr != nil {
		log.Fatalf("Paths Error: %v", pathErr)
	}

	if coverErr := CoverPaths(config, Exec(), paths); coverErr != nil {
		log.Fatalf("Coverage Error: %v", coverErr)
	}
}

func toPaths(args []string) ([]PackagePath, error) {
	paths := make([]PackagePath, 0)
	for _, arg := range args {
		paths = append(paths, PackagePath(arg))
	}
	return paths, nil
}

type PackagePath string
type Package interface {
	FlatName() string
	Name() string
}

type Go interface {
	GetPackages(paths []PackagePath) ([]Package, error)
	CoverageProfile(pack Package, short bool, profilePath string) (bool, error)
	CoverageReport(profilePath string, reportPath string) error
}

func CoverPaths(conf *GoCovConfig, goCmd Go, paths []PackagePath) error {
	packs, pathErr := goCmd.GetPackages(paths)
	if pathErr != nil {
		return pathErr
	}

	for _, pack := range packs {
		if err := CoverPackage(conf, goCmd, pack); err != nil {
			return nil
		}
	}

	return nil
}

func CoverPackage(conf *GoCovConfig, goCmd Go, pack Package) error {
	out := conf.OutputDir.profile(pack.FlatName())
	if !conf.KeepProfile {
		defer os.Remove(out)
	}

	covered, coverErr := goCmd.CoverageProfile(pack, conf.ShortTests, out)
	if coverErr != nil {
		return coverErr
	}

	if covered {
		htmlF := conf.OutputDir.html(pack.FlatName())
		if reportErr := goCmd.CoverageReport(out, htmlF); reportErr != nil {
			fmt.Printf("RPTERR %v", reportErr)
			return reportErr
		}

		if conf.OpenInBrowser {
			fmt.Printf("HTML %v", htmlF)
			browser.OpenFile(htmlF)
		}
	}

	return nil
}
