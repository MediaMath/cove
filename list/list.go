// Package list provides a simplified facade around 'go list'
package list

import (
	"encoding/json"
	"io"

	"github.com/MediaMath/cove/gocmd"
)

// Packages gets all packages that match any of the paths.
// The package list will only contain 1 entry per package, but order is not defined.
// Invalid paths will generate an error, but will not stop the evaluation of the other paths.
func Packages(paths ...string) ([]string, error) {
	return gocmd.Go("list", paths...).StdOutLines()
}

// Json takes a SINGLE fully qualified package import path and decodes the 'go list -json' response.
// See $GOROOT/src/cmd/go/list.go for documentation on the json output.
func Json(pack string, v interface{}) error {
	return gocmd.Go("list", "-json", pack).Receive(func(stdout io.Reader) error {
		return json.NewDecoder(stdout).Decode(v)
	})
}
