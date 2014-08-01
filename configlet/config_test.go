package configlet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrokenConfig(t *testing.T) {
	_, err := Load("./fixtures/broken.json")
	if err == nil {
		t.Errorf("Expected Load() to complain that it couldn't parse the JSON")
	}
}

func TestValidConfig(t *testing.T) {
	path := "./fixtures/valid.json"
	_, err := Load(path)
	if err != nil {
		t.Errorf("Config at %s should be valid, but barfed: %s", path, err)
	}
}

func TestIgnoredDirsIsUnique(t *testing.T) {
	path := "./fixtures/valid.json"
	c, err := Load(path)
	assert.Nil(t, err)

	expected := []string{".git", "bin", "fig", "ignored"}
	actual := c.IgnoredDirs()
	assert.Equal(t, expected, actual)
}
