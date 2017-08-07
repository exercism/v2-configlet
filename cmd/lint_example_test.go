package cmd

import "path/filepath"

func ExampleLint() {
	lintTrack(filepath.FromSlash("../fixtures/numbers"))
	// Output:
	// -> An exercise with slug 'bajillion' is referenced in config.json, but no implementation was found.
	// -> The implementation for 'three' is missing an example solution.
	// -> The implementation for 'two' is missing a test suite.
	// -> An implementation for 'zero' was found, but config.json specifies that it should be foregone (not implemented).
}

func ExampleLintMaintainers() {
	lintTrack(filepath.FromSlash("../fixtures/broken-maintainers"))
	// Output:
	// -> invalid config ../fixtures/broken-maintainers/config/maintainers.json -- invalid character '}' looking for beginning of object key string
}
