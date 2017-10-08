package cmd

// Because the tree cmd concerns itself with writing to stdout the
// more tests are in the example tests. This is concerned with non-output
// related situations.
import (
	"path/filepath"
	"testing"
)

func TestGivenTrackPath(t *testing.T) {
	err := treeTrack(filepath.FromSlash("../fixtures/tree"))

	if err != nil {
		t.Error("should discover config.json given path to directory.")
	}
}

func TestGivenFilename(t *testing.T) {
	err := treeTrack(filepath.FromSlash("../fixtures/tree/config.json"))

	if err != nil {
		t.Error("should open config.json given path to file.")
	}
}

func TestMissingFileError(t *testing.T) {
	err := treeTrack(filepath.FromSlash("../fixtures/tree/non-existing-config.json"))

	if err == nil {
		t.Error("should error for non-existing configuration file.")
	}
}
