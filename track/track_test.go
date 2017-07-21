package track

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTrack(t *testing.T) {
	track, err := New("../fixtures/numbers")
	assert.NoError(t, err)

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
