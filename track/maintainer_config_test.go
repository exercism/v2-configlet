package track

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
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

func TestNoChangeWhenMarshalingAcceptableMaintainerConfig(t *testing.T) {
	filename := "../fixtures/format/formatted/config/maintainers.json"
	src, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		t.Fatal(err)
	}
	mCfg := MaintainerConfig{}
	if err := mCfg.NewConfigFromFile(filename); err != nil {
		t.Fatal(err)
	}
	dst, err := mCfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(src), fmt.Sprintf("%s\n", dst))
}
