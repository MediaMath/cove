package list

import (
	"encoding/json"
	"io"

	"github.com/MediaMath/cove/gocmd"
)

func Packages(paths ...string) ([]string, error) {
	return gocmd.Go("list", paths...).StdOutLines()
}

func Json(pack string, v interface{}) error {
	return gocmd.Go("list", "-json", pack).Receive(func(stdout io.Reader) error {
		return json.NewDecoder(stdout).Decode(v)
	})
}
