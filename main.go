package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/MediaMath/gocov/packages"
	"github.com/MediaMath/gocov/relative"
	"github.com/pkg/browser"
)

func main() {
	wd, wdErr := os.Getwd()

	if wdErr != nil {
		log.Fatal(wdErr)
	}

	outputFiles, outputErr := existingDirOutputFiles(filepath.Join(wd, "tmp", "coverage2"))
	if outputErr != nil {
		log.Fatal(outputErr)
	}

	conf := config{true, true, outputFiles, false, false}
	err := gocovForAll(conf, wd, []string{wd})
	if err != nil {
		log.Fatal(err)
	}

}

type config struct {
	recurse     bool
	browser     bool
	output      outputFiles
	short       bool
	keepProfile bool
}

func gocovForAll(conf config, wd string, paths []string) error {
	packs, getErr := packages.GetAll(conf.recurse, paths...)
	if getErr != nil {
		return getErr
	}

	for _, pack := range packs {
		gocovForPack(conf, wd, pack)
	}

	return nil
}

func gocovForPack(conf config, wd string, pack string) error {
	rel := relative.To(pack, wd)

	out := conf.output.profile(rel.FlatName())
	if !conf.keepProfile {
		defer os.Remove(out)
	}

	covered, coverErr := buildProfile(fmt.Sprintf("./%v", rel.Path()), out)
	if coverErr != nil {
		return coverErr
	}

	if covered {
		htmlF := conf.output.html(rel.FlatName())
		buildReport(out, htmlF)
		if conf.browser {
			browser.OpenFile(htmlF)
		}
	}

	return nil
}
