package layout

import (
	"os"
	"strings"
)

// hasPackage reports whether a directory is a real Go package. It is a
// package-level seam, defaulting to the filesystem implementation, so tests
// drive the existence decision with a fake and the analysis stays hermetic
// (no direct OS reach-through, per the gomatic dependency-injection standard).
var hasPackage = osHasPackage

// osHasPackage reports whether dir contains at least one non-test Go source
// file, i.e. is a real Go package rather than a bare or empty directory. A
// counterpart directory with no Go source does not satisfy the layout, so
// existence alone is not enough.
func osHasPackage(dir pkgDir) bool {
	entries, err := os.ReadDir(string(dir))
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && isGoSource(fileName(entry.Name())) {
			return true
		}
	}
	return false
}

// fileName is a bare file name within a package directory.
type fileName string

// isGoSource reports whether name is a non-test Go source filename.
func isGoSource(name fileName) bool {
	return strings.HasSuffix(string(name), ".go") && !strings.HasSuffix(string(name), "_test.go")
}
