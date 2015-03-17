package relative

import (
	"path/filepath"
	"strings"
)

type Relative interface {
	Path() string
	FlatName() string
}

func To(path string, relativeTo string) Relative {
	return &relative{path, relativeTo}
}

type relative struct {
	fullPath   string
	relativeTo string
}

func (r *relative) Path() string {
	rel, _ := filepath.Rel(r.relativeTo, r.fullPath)
	return rel
}

func (r *relative) FlatName() string {
	return strings.Replace(r.Path(), "/", "-", -1)
}
