package layout

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis"
)

// parsePass builds a single-file analysis pass for filename and src, capturing
// every diagnostic the analyzer reports into the returned slice pointer.
func parsePass(t *testing.T, filename, src string) (*analysis.Pass, *[]string) {
	t.Helper()
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, src, 0)
	require.NoError(t, err)
	var messages []string
	pass := &analysis.Pass{
		Fset:   fset,
		Files:  []*ast.File{file},
		Report: func(diag analysis.Diagnostic) { messages = append(messages, diag.Message) },
	}
	return pass, &messages
}

// TestPackageDirGuardsPathWithoutSeparator covers the no-separator guard: a file
// name carrying no "/" must be returned unchanged rather than slicing past -1.
func TestPackageDirGuardsPathWithoutSeparator(t *testing.T) {
	pass, _ := parsePass(t, "noslash", "package p")
	assert.Equal(t, "noslash", packageDir(pass))
}

// TestRunResolvesCounterpartThroughSeam proves the injected hasPackage seam
// drives the correspondence decision: a present counterpart yields no
// diagnostic, a missing one yields exactly the counterpart-missing report.
func TestRunResolvesCounterpartThroughSeam(t *testing.T) {
	original := hasPackage
	t.Cleanup(func() { hasPackage = original })

	pass, messages := parsePass(t, "/m/internal/app/commands/greet/command.go", "package greet")

	hasPackage = func(string) bool { return true }
	_, err := run(pass)
	require.NoError(t, err)
	assert.Empty(t, *messages)

	hasPackage = func(string) bool { return false }
	_, err = run(pass)
	require.NoError(t, err)
	assert.Len(t, *messages, 1)
}
