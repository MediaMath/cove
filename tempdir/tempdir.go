package tempdir

import (
	"io/ioutil"
	"os"
)

func In(action func(string)) error {
	f, testdir, err := Make("", getCallerLabel())
	defer f()

	if err != nil {
		return err
	} else {
		action(testdir)
		return nil
	}
}

func Make(path string, prefix string) (func(), string, error) {
	testdir, err := ioutil.TempDir(path, prefix)
	if err != nil {
		return func() {}, "", err
	} else {
		return func() { os.RemoveAll(testdir) }, testdir, nil
	}
}
