package cmd

import (
	"os"
	"path/filepath"

	"github.com/exercism/configlet/ui"
)

func ExampleLint() {
	ui.Out = os.Stdout
	originalNoHTTP := noHTTP
	noHTTP = true
	defer func() { noHTTP = originalNoHTTP }()

	lintTrack(filepath.FromSlash("../fixtures/numbers"))
	// Output:
	// -> An exercise with slug 'bajillion' is referenced in config.json, but no implementation was found.
	// -> The implementation for 'three' is missing an example solution.
	// -> The implementation for 'two' is missing a test suite.
	// -> The exercise 'one' was found in config.json, but does not have a UUID.
	// -> An implementation for 'zero' was found, but config.json specifies that it should be foregone (not implemented).
}

func ExampleLintMaintainers() {
	ui.ErrOut = os.Stdout
	originalNoHTTP := noHTTP
	noHTTP = true
	defer func() { noHTTP = originalNoHTTP }()

	lintTrack(filepath.FromSlash("../fixtures/broken-maintainers"))
	// Output:
	// -> invalid config ../fixtures/broken-maintainers/config/maintainers.json -- invalid character '}' looking for beginning of object key string
}
