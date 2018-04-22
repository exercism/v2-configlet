package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
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
	expectedConfig = strings.TrimSuffix(string(cfg), "\n")

	maintainers, err := ioutil.ReadFile(filepath.FromSlash("../fixtures/format/formatted/maintainers.json"))
	if err != nil {
		log.Fatal(err)
	}
	expectedMaintainers = strings.TrimSuffix(string(maintainers), "\n")
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
		_, actualConfig, err := formatFile(filepath.FromSlash(f), formatTopics, orderConfig)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, expectedConfig, string(actualConfig))
	}
}

var maintainersFiles = []string{
	"../fixtures/format/formatted/maintainers.json",
	"../fixtures/format/malformed/maintainers.json",
	"../fixtures/format/minimised/maintainers.json",
}

func TestMaintainers(t *testing.T) {
	for _, f := range maintainersFiles {
		_, actualMaintainers, err := formatFile(filepath.FromSlash(f), nil, nil)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, expectedMaintainers, string(actualMaintainers))
	}
}
