package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFmtCommand(t *testing.T) {
	trackCfg, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/config.json"))
	if err != nil {
		t.Fatal(err)
	}

	maintainerCfg, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/config/maintainers.json"))
	if err != nil {
		t.Fatal(err)
	}

	// The fmt command does not rewrite a correctly formatted file.
	formattedDir, err := ioutil.TempDir("", "formatted")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(formattedDir)

	fmtTest = false
	fmtVerbose = false
	runFmt("../fixtures/format/formatted/", formattedDir)

	_, err = os.Stat(filepath.Join(formattedDir, "config.json"))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(formattedDir, "config", "maintainers.json"))
	assert.True(t, os.IsNotExist(err))

	// It rewrites an incorrectly formatted file.
	malformedDir, err := ioutil.TempDir("", "malformed")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(malformedDir)

	fmtTest = false
	fmtVerbose = false
	runFmt("../fixtures/format/malformed/", malformedDir)

	track, err := ioutil.ReadFile(filepath.Join(malformedDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, track, trackCfg)

	maintainer, err := ioutil.ReadFile(filepath.Join(malformedDir, "config", "maintainers.json"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, maintainer, maintainerCfg)

	// It does not rewrite an incorrectly formatted file when the diff opton is true
	unformattedDir, err := ioutil.TempDir("", "unformatted")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(formattedDir)

	fmtTest = true
	fmtVerbose = false
	runFmt("../fixtures/format/unformatted/", unformattedDir)

	_, err = os.Stat(filepath.Join(unformattedDir, "config.json"))
	assert.True(t, os.IsNotExist(err))

	_, err = os.Stat(filepath.Join(unformattedDir, "config", "maintainers.json"))
	assert.True(t, os.IsNotExist(err))

}

func TestSemanticsOfMissingTopics(t *testing.T) {
	semanticsDir, err := ioutil.TempDir("", "semantics")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(semanticsDir)

	fmtTest = false
	fmtVerbose = false
	runFmt("../fixtures/format/semantics/", semanticsDir)

	// No change; nothing should be written to out dir.
	_, err = os.Stat(filepath.Join(semanticsDir, "config.json"))
	assert.True(t, os.IsNotExist(err))
}
