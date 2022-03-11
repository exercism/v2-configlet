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
	// -> The exercise with slug 'not-implemented' is not implemented in this track.
	// -> The exercise with slug 'old-but-gold' exists in this track, but has been removed from the specifications repository.
	// -> The exercise with slug 'old-but-present' exists in this track, but has been deprecated in the specifications repository.
}
