package track

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBrokenMaintainerConfig(t *testing.T) {
	if _, err := NewMaintainerConfig("../fixtures/broken-maintainers/config/maintainers.json"); err == nil {
		t.Errorf("Expected broken JSON")
	}
}

func TestValidMaintainerConfig(t *testing.T) {
	mc, err := NewMaintainerConfig("../fixtures/numbers/config/maintainers.json")
	if err != nil {
		t.Errorf("Expected valid JSON: %s", err)
	}
	assert.Equal(t, "alice", mc.Maintainers[0].Username)
}

func TestIgnoreMissingMaintainerConfig(t *testing.T) {
	_, err := NewMaintainerConfig("../fixtures/no-such-file.json")
	assert.NoError(t, err)
}
