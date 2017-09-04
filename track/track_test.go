package track

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrack(t *testing.T) {
	track, err := New("../fixtures/numbers")
	assert.NoError(t, err)

	assert.Equal(t, "numbers", track.ID)
	assert.Equal(t, "Numbers", track.Config.Language)

	slugs := map[string]bool{
		"zero":  false,
		"one":   false,
		"two":   false,
		"three": false,
	}

	assert.Equal(t, len(slugs), len(track.Exercises))

	for _, exercise := range track.Exercises {
		slugs[exercise.Slug] = true
	}

	for slug, ok := range slugs {
		if !ok {
			t.Errorf("Expected to find exercise %s", slug)
		}
	}
}

func TestTrackID(t *testing.T) {

	tests := []struct {
		root     string
		path     string
		expected string
	}{
		{"../fixtures", "numbers", "numbers"},
		{"../fixtures/numbers", ".", "numbers"},
	}

	cwd, _ := os.Getwd()
	defer func() { os.Chdir(cwd) }()
	for _, test := range tests {
		err := os.Chdir(test.root)
		assert.NoError(t, err)

		track, err := New(test.path)
		assert.NoError(t, err)

		assert.Equal(t, test.expected, track.ID)

		// reset working directory for each test
		os.Chdir(cwd)
	}
}
