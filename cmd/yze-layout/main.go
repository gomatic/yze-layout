// Command yze-layout runs the layout analyzer as a standalone go/analysis
// checker (text and -json output, and as a `go vet -vettool`).
package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	layout "github.com/gomatic/yze-layout"
)

// run is the analysis entry point, indirected so the binary's wiring is testable
// without invoking the real driver (which loads packages and exits the process).
var run = singlechecker.Main

func main() { run(layout.Analyzer) }
