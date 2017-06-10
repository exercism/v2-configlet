package configlet

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigHasDefaultSolutionPattern(t *testing.T) {
	c := NewConfig()
	expected := "[Ee]xample"
	assert.Equal(t, expected, c.SolutionPattern)
}

func TestBrokenConfig(t *testing.T) {
	if _, err := Load("./fixtures/broken.json"); err == nil {
		t.Errorf("Expected Load() to complain that it couldn't parse the JSON")
	}
}

func TestValidConfig(t *testing.T) {
	if _, err := Load("./fixtures/valid.json"); err != nil {
		t.Errorf("Expected valid.json to contain valid JSON: %s", err)
	}
}

func TestConfigSlugs(t *testing.T) {
	expectedSlugs := []string{
		"apple",
		"banana",
		"cherimoya",
	}

	path := "./fixtures/exercises.json"
	c, err := Load(path)
	if err != nil {
		t.Errorf("failed to load config at %s.", path)
	}

	actualSlugs := c.Slugs()
	if len(actualSlugs) != len(expectedSlugs) {
		t.Errorf("%s: got %d slugs, want %d", path, len(actualSlugs), len(expectedSlugs))
	}

	for i, slug := range c.Slugs() {
		if expectedSlugs[i] != slug {
			t.Errorf("%s - slugs[%d]: expected '%s', got '%s'", path, i, expectedSlugs[i], slug)
		}
	}
}

func TestIgnoredDirsIsUnique(t *testing.T) {
	path := "./fixtures/valid.json"
	c, err := Load(path)
	assert.Nil(t, err)

	expected := []string{"bin", "fig", "ignored", "img"}
	actual := c.IgnoredDirs()

	sort.Strings(expected)
	sort.Strings(actual)

	assert.Equal(t, expected, actual)
}
