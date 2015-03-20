package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/MediaMath/cove"
	"github.com/MediaMath/cove/cmd"
)

//gosh - Get Over SsH - gets go packages via specific uri's.
//the arguments are of the form <package_name>,<uri> - github.com/MediaMath/cove,git@github.com:MediaMath/cove
func main() {
	overwrite := flag.Bool("f", false, "Overwrite existing packages")
	flag.Parse()

	goshMap, parseErr := getMap(flag.Args())
	if parseErr != nil {
		log.Fatal(parseErr)
	}

	if gosh(*overwrite, goshMap) {
		os.Exit(1)
	}
}

func getMap(args []string) (map[string]string, error) {
	goshMap := make(map[string]string)
	for _, arg := range args {
		pair := strings.Split(arg, ",")
		if len(pair) != 2 {
			return nil, fmt.Errorf("Arguments are unparseable: %v", strings.Join(args, " "))
		}

		goshMap[pair[0]] = pair[1]
	}

	return goshMap, nil
}

func gosh(overwrite bool, goshMap map[string]string) bool {
	hadError := false
	for pack, location := range goshMap {

		if overwrite || !cove.PackageExists(pack) {
			if logError(clone(location, to(pack))) {
				hadError = true
				continue
			}
		}

		dependencies, deppErr := cove.MissingDependencies(pack)
		if logError(deppErr) {
			hadError = true
			continue
		}

		if getDependencies(goshMap, dependencies) {
			hadError = true
		}

	}

	return hadError
}

func getDependencies(goshMap map[string]string, dependencies []string) bool {
	hadError := false
	for _, dep := range dependencies {
		if _, inGoshMap := goshMap[dep]; inGoshMap {
			continue
		}

		fmt.Printf("go get %v\n", dep)
		if logError(cove.Get(dep)) {
			hadError = true
		}
	}

	return hadError
}

func logError(err error) bool {
	if err != nil {
		log.Printf("%s\n", err)
		return true
	}

	return false
}

func clone(src string, destination string) error {
	vcsCopy := func(copyTemp string) error { return vcsClone(src, copyTemp) }
	return overwrite(destination, vcsCopy)
}

func vcsClone(src string, destination string) error {
	fmt.Printf("git clone %v\n", src)
	return cmd.Run(exec.Command("git", "clone", "--depth=1", src, destination))
}

func to(pack string) string {
	return filepath.Join(os.Getenv("GOPATH"), "src", pack)
}

func overwrite(destination string, copyTo func(string) error) error {
	copyTemp, copyTempErr := ioutil.TempDir("", "overwrite-copy")
	if copyTempErr != nil {
		return copyTempErr
	}
	defer os.RemoveAll(copyTemp)

	if copyErr := copyTo(copyTemp); copyErr != nil {
		return copyErr
	}

	oldDestinationTemp, oldTempErr := ioutil.TempDir("", "overwrite-old")
	if oldTempErr != nil {
		return oldTempErr
	}

	//if destination does not already exist create it
	os.MkdirAll(destination, 0744)

	if renameErr := os.Rename(destination, oldDestinationTemp); renameErr != nil {
		return renameErr
	}

	if replaceErr := os.Rename(copyTemp, destination); replaceErr != nil {
		if undoErr := os.Rename(oldDestinationTemp, destination); undoErr != nil {
			return fmt.Errorf("%v\n%v", replaceErr, undoErr)
		}

		return replaceErr
	}

	os.RemoveAll(oldDestinationTemp)
	return nil
}
