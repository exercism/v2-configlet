package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	expectedConfig      string
	expectedMaintainers string
)

func TestMain(m *testing.M) {
	cfg, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/config.json"))
	if err != nil {
		log.Fatal(err)
	}
	expectedConfig = string(cfg)

	maintainers, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/config/maintainers.json"))
	if err != nil {
		log.Fatal(err)
	}
	expectedMaintainers = string(maintainers)
	result := m.Run()
	os.Exit(result)
}

var configFiles = []string{
	"../fixtures/format/formatted/config.json",
	"../fixtures/format/malformed/config.json",
	"../fixtures/format/minimised/config.json",
}

func TestFormat(t *testing.T) {
	for _, f := range configFiles {
		tmp, err := ioutil.TempFile(os.TempDir(), "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmp.Name())

		_, _, err = formatFile(filepath.FromSlash(f), tmp.Name(), formatTopics, orderConfig)
		if err != nil {
			log.Fatal(err)
		}

		actualConfig, err := ioutil.ReadFile(tmp.Name())
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, expectedConfig, string(actualConfig))
	}
}

var maintainersFiles = []string{
	"../fixtures/format/formatted/config/maintainers.json",
	"../fixtures/format/malformed/config/maintainers.json",
	"../fixtures/format/minimised/config/maintainers.json",
}

func TestMaintainers(t *testing.T) {
	for _, f := range maintainersFiles {
		tmp, err := ioutil.TempFile(os.TempDir(), "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmp.Name())

		_, _, err = formatFile(filepath.FromSlash(f), tmp.Name(), nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		actualMaintainers, err := ioutil.ReadFile(tmp.Name())
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, expectedMaintainers, string(actualMaintainers))
	}
}

func TestNoChangeOnFormattingCompliantConfig(t *testing.T) {
	filename := "../fixtures/format/formatted/config.json"

	tmp, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	// Read original from source.
	src, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		t.Fatal(err)
	}

	// Run it through the formatter.
	_, dst, err := formatFile(filepath.FromSlash(filename), tmp.Name(), formatTopics, orderConfig)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(src), string(dst))
}

func TestSemanticsOfMissingTopics(t *testing.T) {
	f := "../fixtures/format/semantics/config.json"

	tmp, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())

	// Read original from source
	src, err := ioutil.ReadFile(filepath.FromSlash(f))
	if err != nil {
		t.Fatal(err)
	}

	_, dst, err := formatFile(filepath.FromSlash(f), tmp.Name(), formatTopics, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(src), string(dst))
}
