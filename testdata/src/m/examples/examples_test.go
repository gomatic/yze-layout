package examples_test

// examples_test.go is the sole Go file in its directory and is an external test
// file, so the driver delivers a base-package pass with no syntax files. The
// directory is not a three-tier command or domain package, so the analyzer must
// run clean (no diagnostics) instead of indexing pass.Files[0]. Regression for
// the examples-only crash.
func ExampleLayout() {}
