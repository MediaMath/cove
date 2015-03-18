package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Exec() Go {
	return &e{}
}

type e struct{}

func getShort(short bool) string {
	if short {
		return "-short"
	} else {
		return ""
	}
}

func (a *e) GetPackages(paths []PackagePath) ([]Package, error) {

	packs := make([]Package, 0)
	for _, path := range paths {
		next, err := getPackagesForPath(path)
		if err != nil {
			return packs, err
		}
		packs = append(packs, next...)
	}

	return packs, nil
}

func getPackagesForPath(path PackagePath) ([]Package, error) {
	output, err := exec.Command("go", "list", "-e", string(path)).Output()
	outString := ""
	if err == nil {
		outString = fmt.Sprintf("%s", output)
	}

	packNames := strings.Split(outString, "\n")
	packs := make([]Package, 0)
	for _, p := range packNames {
		packs = append(packs, pack(strings.TrimSpace(p)))
	}

	return packs, nil
}

type pack string

func (p pack) Name() string {
	return string(p)
}

func (p pack) FlatName() string {
	return strings.Replace(p.Name(), "/", ".", -1)
}

func (a *e) CoverageProfile(pack Package, short bool, outFile string) (bool, error) {
	command := []string{"test", pack.Name(), fmt.Sprintf("-coverprofile=%s", outFile)}
	if short {
		command = append(command, "-short")
	}

	cmd := exec.Command("go", command...)
	cmd.Output()

	_, err := os.Stat(outFile)
	return err == nil, nil
}

func (a *e) CoverageReport(profilePath string, reportPath string) error {

	cmd := exec.Command("go", "tool", "cover", fmt.Sprintf("-html=%s", profilePath), "-o", reportPath)
	_, err := cmd.Output()
	return err
}
