package packages

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/MediaMath/gocov/tempdir"
)

func TestNonExistentPathIsNotAPackage(t *testing.T) {
	madeUp := filepath.Join("madeup", "foo", "bar", "moo")
	if IsPackage(madeUp) {
		t.Errorf("made up: %v should not be a package", madeUp)
	}
}

func TestDirThatExistsWithNoFilesInItIsNotAPackage(t *testing.T) {
	tempdir.In(func(dir string) {
		if IsPackage(dir) {
			t.Errorf("empty dir: %v should not be a package", dir)
		}
	})
}

func TestFileThatExistsIsNotAPackage(t *testing.T) {
	tempdir.In(func(dir string) {
		fileName := makeTestPackage(dir)
		if IsPackage(fileName) {
			t.Errorf("file: %v should not be a package", fileName)
		}
	})
}

func TestDirWithFileExistsButNoGoFilesIsNotAPackage(t *testing.T) {
	tempdir.In(func(dir string) {
		makeTestFile(dir, "file.gog")
		if IsPackage(dir) {
			t.Errorf("dir with no go files: %v should not be a package", dir)
		}
	})
}

func TestDirWithGoFileIsPackage(t *testing.T) {
	tempdir.In(func(dir string) {
		makeTestPackage(dir)
		if !IsPackage(dir) {
			t.Errorf("dir with go file: %v should be a package", dir)
		}
	})
}

func TestContainsIsEmptyWithNoPackagesInRoot(t *testing.T) {
	tempdir.In(func(dir string) {
		packs, _ := GetAll(true, dir)
		if len(packs) > 0 {
			t.Errorf("Should have been empty: %v", packs)
		}
	})
}

func TestGetContainsRootIfItIsAPackage(t *testing.T) {
	tempdir.In(func(dir string) {
		makeTestPackage(dir)
		packs, _ := GetAll(true, dir)
		if !contains(packs, dir) {
			t.Errorf("%v did not contain %v", packs, dir)
		}
	})
}

func TestGetRecurses(t *testing.T) {
	tempdir.In(func(dir string) {
		f1, pack1, _ := tempdir.Make(dir, "pack1")
		defer f1()
		makeTestPackage(pack1)

		f2, nopack, _ := tempdir.Make(dir, "nopack")
		defer f2()

		f3, pack2, _ := tempdir.Make(pack1, "pack2")
		defer f3()
		makeTestPackage(pack2)

		packs, _ := GetAll(true, dir)
		if !contains(packs, pack1) {
			t.Errorf("%v did not contain %v", packs, pack1)
		}

		if !contains(packs, pack2) {
			t.Errorf("%v did not contain %v", packs, pack2)
		}

		if contains(packs, nopack) {
			t.Errorf("%v did contained %v", packs, nopack)
		}
	})
}

func TestGetSiblingRecurses(t *testing.T) {
	tempdir.In(func(dir string) {
		f1, pack1, _ := tempdir.Make(dir, "pack1")
		defer f1()
		makeTestPackage(pack1)

		f2, nopack1, _ := tempdir.Make(dir, "nopack1")
		defer f2()

		f3, pack2, _ := tempdir.Make(pack1, "pack2")
		defer f3()
		makeTestPackage(pack2)

		f4, pack3, _ := tempdir.Make(dir, "pack3")
		defer f4()
		makeTestPackage(pack3)

		f5, pack4, _ := tempdir.Make(pack3, "pack4")
		defer f5()
		makeTestPackage(pack4)

		f6, nopack2, _ := tempdir.Make(pack3, "nopack2")
		defer f6()

		packs, _ := GetAll(true, pack1, nopack1, pack3)
		if !contains(packs, pack1) {
			t.Errorf("%v did not contain %v", packs, pack1)
		}

		if !contains(packs, pack2) {
			t.Errorf("%v did not contain %v", packs, pack2)
		}

		if !contains(packs, pack3) {
			t.Errorf("%v did not contain %v", packs, pack3)
		}
		if !contains(packs, pack4) {
			t.Errorf("%v did not contain %v", packs, pack4)
		}

		if contains(packs, nopack1) {
			t.Errorf("%v did contained %v", packs, nopack1)
		}

		if contains(packs, nopack2) {
			t.Errorf("%v did contained %v", packs, nopack2)
		}
	})
}

func TestGetSiblingNoRecurses(t *testing.T) {
	tempdir.In(func(dir string) {
		f1, pack1, _ := tempdir.Make(dir, "pack1")
		defer f1()
		makeTestPackage(pack1)

		f2, nopack1, _ := tempdir.Make(dir, "nopack1")
		defer f2()

		f3, pack2, _ := tempdir.Make(pack1, "pack2")
		defer f3()
		makeTestPackage(pack2)

		f4, pack3, _ := tempdir.Make(dir, "pack3")
		defer f4()
		makeTestPackage(pack3)

		f5, pack4, _ := tempdir.Make(pack3, "pack4")
		defer f5()
		makeTestPackage(pack4)

		f6, nopack2, _ := tempdir.Make(pack3, "nopack2")
		defer f6()

		packs, _ := GetAll(false, pack1, nopack1, pack3)
		if !contains(packs, pack1) {
			t.Errorf("%v did not contain %v", packs, pack1)
		}

		if contains(packs, pack2) {
			t.Errorf("%v did contained %v", packs, pack2)
		}

		if !contains(packs, pack3) {
			t.Errorf("%v did not contain %v", packs, pack3)
		}
		if contains(packs, pack4) {
			t.Errorf("%v did contained %v", packs, pack4)
		}

		if contains(packs, nopack1) {
			t.Errorf("%v did contained %v", packs, nopack1)
		}

		if contains(packs, nopack2) {
			t.Errorf("%v did contained %v", packs, nopack2)
		}
	})
}
func TestGetDoesntRecurse(t *testing.T) {
	tempdir.In(func(dir string) {
		f1, pack1, _ := tempdir.Make(dir, "pack1")
		defer f1()
		makeTestPackage(pack1)

		f2, _, _ := tempdir.Make(dir, "nopack")
		defer f2()

		f3, pack2, _ := tempdir.Make(pack1, "pack2")
		defer f3()
		makeTestPackage(pack2)

		packs, _ := GetAll(false, dir)
		if contains(packs, pack1) {
			t.Errorf("%v did not contain %v", packs, pack1)
		}
	})
}

func contains(packs []string, pack string) bool {
	for _, p := range packs {
		if p == pack {
			return true
		}
	}

	return false
}

func makeTestPackage(path string) string {
	return makeTestFile(path, "file.go")
}

func makeTestFile(path string, name string) string {
	fileName := filepath.Join(path, name)
	ioutil.WriteFile(fileName, []byte("//test"), 0644)
	return fileName
}
