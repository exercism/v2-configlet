package track

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrokenConfig(t *testing.T) {
	if _, err := NewConfig("../fixtures/broken.json"); err == nil {
		t.Errorf("Expected broken JSON")
	}
}

func TestValidConfig(t *testing.T) {
	c, err := NewConfig("../fixtures/numbers/config.json")
	if err != nil {
		t.Errorf("Expected valid JSON: %s", err)
	}
	assert.Equal(t, "Numbers", c.Language)
}

func TestDefaultSolutionPattern(t *testing.T) {
	c, err := NewConfig("../fixtures/empty.json")
	assert.NoError(t, err)
	assert.Equal(t, "[Ee]xample", c.SolutionPattern)
}

func TestDefaultTestPattern(t *testing.T) {
	c, err := NewConfig("../fixtures/empty.json")
	assert.NoError(t, err)
	assert.Equal(t, "(?i)test", c.TestPattern)
}
