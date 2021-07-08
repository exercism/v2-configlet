package cmd

import (
	"os"
	"path/filepath"

	"github.com/exercism/configlet/ui"
)

func ExampleImplementationStatus() {
	ui.Out = os.Stdout

	runImplementationStatus(
		filepath.FromSlash("../fixtures/implementation-status/problem-specifications"),
		filepath.FromSlash("../fixtures/implementation-status/track"))
	// Output:
	// -> The exercise 'not-implemented' does not exist in this track.
	// -> The exercise 'old-but-present' exists in this track but is deprecated.
	// -> The exercise 'old-but-gold' exists in the track but is not in the problem-specifications repository.
}
