package cmd

// Because the visualizer concerns itself with writing to stdout the
// more tests are in the example tests. This is concerned with non-output
// related situations.
import (
	"path/filepath"
	"testing"
)

func TestNoFileError(t *testing.T) {
	err := visualizeTrack(filepath.FromSlash("../fixtures/visualize/non-existing-config.json"))

	if err == nil {
		t.Error("expected error for non-existing configuration file")
	}
}
