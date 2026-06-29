// Package layout provides a go/analysis analyzer enforcing the cross-package
// correspondence of the gomatic three-tier CLI layout: every
// internal/app/commands/<cmd> package has a matching internal/domain/<cmd>
// package, and vice versa. Each package checks its own counterpart on the
// filesystem, so the two directions are reported without duplication.
package layout

import (
	"os"
	"strings"

	goyze "github.com/gomatic/go-yze"
	"golang.org/x/tools/go/analysis"
)

const (
	commandSegment = "/internal/app/commands/"
	domainSegment  = "/internal/domain/"
)

// Analyzer reports command or domain packages whose counterpart is missing.
var Analyzer = &analysis.Analyzer{
	Name: "layout",
	Doc:  "reports three-tier packages whose corresponding command or domain package is missing",
	Run:  run,
}

// Registration declares this analyzer to the yze framework.
var Registration = goyze.Registration{
	Name:       "layout",
	Categories: []goyze.Category{"structure"},
	URL:        "https://docs.gomatic.dev/yze/layout",
	Analyzer:   Analyzer,
}

// hasPackage reports whether a directory is a real Go package. It is a
// package-level seam, defaulting to the filesystem implementation, so tests
// drive the existence decision with a fake and the analysis stays hermetic
// (no direct OS reach-through, per the gomatic dependency-injection standard).
var hasPackage = osHasPackage

// run reports when a command or domain package has no counterpart package.
//
// A package whose only Go files are external test files (an examples-only
// directory) is delivered by the driver as a base-package pass with no syntax
// files. There is no file to locate the package directory or anchor a report
// on, and such a directory is never a three-tier command or domain package, so
// the pass is a no-op rather than an index-out-of-range panic on pass.Files[0].
func run(pass *analysis.Pass) (any, error) {
	if len(pass.Files) == 0 {
		return nil, nil
	}
	counterpart, message, ok := counterpartOf(packageDir(pass))
	if ok && !hasPackage(counterpart) {
		pass.Reportf(pass.Files[0].Name.Pos(), "%s", message)
	}
	return nil, nil
}

// packageDir returns the filesystem directory of the analyzed package.
func packageDir(pass *analysis.Pass) string {
	name := pass.Fset.File(pass.Files[0].Pos()).Name()
	idx := strings.LastIndex(name, "/")
	if idx < 0 {
		return name
	}
	return name[:idx]
}

// counterpartOf returns the directory that must exist for a command or domain
// package, the diagnostic to emit if it is missing, and whether dir is a
// three-tier package at all.
func counterpartOf(dir string) (string, string, bool) {
	if strings.Contains(dir, commandSegment) {
		return strings.Replace(dir, commandSegment, domainSegment, 1),
			"command package has no corresponding internal/domain package", true
	}
	if strings.Contains(dir, domainSegment) {
		return strings.Replace(dir, domainSegment, commandSegment, 1),
			"domain package has no corresponding internal/app/commands package", true
	}
	return "", "", false
}

// osHasPackage reports whether dir contains at least one non-test Go source
// file, i.e. is a real Go package rather than a bare or empty directory. A
// counterpart directory with no Go source does not satisfy the layout, so
// existence alone is not enough.
func osHasPackage(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() && isGoSource(entry.Name()) {
			return true
		}
	}
	return false
}

// isGoSource reports whether name is a non-test Go source filename.
func isGoSource(name string) bool {
	return strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go")
}
