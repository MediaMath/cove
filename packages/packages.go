package packages

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func IsPackage(path string) bool {
	paths, readErr := ioutil.ReadDir(path)
	if readErr != nil {
		return false
	}

	for _, path := range paths {
		if isGoFile(path) {
			return true
		}
	}

	return false
}

func isGoFile(file os.FileInfo) bool {
	return filepath.Ext(file.Name()) == ".go"
}

func walk(root string, visitor func(string, error) error) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if IsPackage(path) {
			err := visitor(path, nil)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func GetAll(recurse bool, paths ...string) ([]string, error) {
	if recurse {
		return flatMapRecursionPacks(paths)
	} else {
		return get(paths...), nil
	}
}

func flatMapRecursionPacks(paths []string) ([]string, error) {
	flatMapped := make([]string, 0)
	for _, root := range paths {
		recursed, err := getRecurse(root)
		if err != nil {
			return flatMapped, err
		}
		flatMapped = append(flatMapped, recursed...)
	}

	return flatMapped, nil
}

func get(paths ...string) []string {
	packs := make([]string, 0)
	for _, path := range paths {
		if IsPackage(path) {
			packs = append(packs, path)
		}
	}

	return packs
}

func getRecurse(root string) ([]string, error) {
	collected := make([]string, 0)
	walkErr := walk(root, func(pack string, err error) error {
		collected = append(collected, pack)
		return err
	})

	return collected, walkErr
}
