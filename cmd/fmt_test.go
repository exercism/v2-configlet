package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoChangeOnFormattingCompliantConfig(t *testing.T) {
	filename := "../fixtures/format/formatted/config.json"

	tmp, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	// Run it through the formatter.
	if _, err := formatFile(filepath.FromSlash(filename), tmp.Name(), formatTopics, orderConfig); err != nil {
		t.Fatal(err)
	}

	// No change; nothing should be written to outfile.
	dst, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(dst))
}

func TestSemanticsOfMissingTopics(t *testing.T) {
	f := "../fixtures/format/semantics/config.json"

	tmp, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	if _, err := formatFile(filepath.FromSlash(f), tmp.Name(), formatTopics, nil); err != nil {
		t.Fatal(err)
	}

	// No change; nothing should be written to outfile.
	dst, err := ioutil.ReadFile(tmp.Name())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", string(dst))
}

func TestFmtCommand(t *testing.T) {
	cfgTrack, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/config.json"))
	if err != nil {
		t.Fatal(err)
	}

	cfgMaintainer, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/config/maintainers.json"))
	if err != nil {
		t.Fatal(err)
	}

	// The fmt command does not rewrite a correctly formatted file.
	formattedDir, err := ioutil.TempDir("", "formatted")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(formattedDir)
	runFmt("../fixtures/format/formatted/", formattedDir, false)

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
	runFmt("../fixtures/format/malformed/", malformedDir, false)

	track, err := ioutil.ReadFile(filepath.Join(malformedDir, "config.json"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, track, cfgTrack)

	maintainer, err := ioutil.ReadFile(filepath.Join(malformedDir, "config", "maintainers.json"))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, maintainer, cfgMaintainer)
}
