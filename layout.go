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
	Group:      "go",
	Categories: []goyze.Category{"structure"},
	URL:        "https://docs.gomatic.dev/yze/go/layout",
	Analyzer:   Analyzer,
}

// run reports when a command or domain package has no counterpart directory.
func run(pass *analysis.Pass) (any, error) {
	counterpart, message, ok := counterpartOf(packageDir(pass))
	if ok && !isDir(counterpart) {
		pass.Reportf(pass.Files[0].Name.Pos(), "%s", message)
	}
	return nil, nil
}

// packageDir returns the filesystem directory of the analyzed package.
func packageDir(pass *analysis.Pass) string {
	first := pass.Fset.File(pass.Files[0].Pos()).Name()
	return first[:strings.LastIndex(first, "/")]
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

// isDir reports whether path is an existing directory.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
