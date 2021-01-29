package track

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
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

func TestDefaultIgnorePattern(t *testing.T) {
	c, err := NewConfig("../fixtures/empty.json")
	assert.NoError(t, err)
	assert.Equal(t, "[Ee]xample", c.IgnorePattern)
}

func TestNoChangeWhenMarshalingAcceptableConfig(t *testing.T) {
	filename := "../fixtures/format/formatted/config.json"
	src, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		t.Fatal(err)
	}

	cfg := Config{}
	if err := cfg.LoadFromFile(filename); err != nil {
		t.Fatal(err)
	}
	dst, err := cfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, string(src), fmt.Sprintf("%s\n", dst))
}

func TestMarshalingSortsTopics(t *testing.T) {
	src := `{"exercises": [{"topics": ["banana","cherry","apple"]}]}`

	var srcCfg Config
	if err := json.NewDecoder(strings.NewReader(src)).Decode(&srcCfg); err != nil {
		t.Fatal(err)
	}

	dst, err := srcCfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	var dstCfg Config
	if err := json.NewDecoder(bytes.NewReader(dst)).Decode(&dstCfg); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"apple", "banana", "cherry"}, dstCfg.Exercises[0].Topics)
}

func TestMarshalingNormalizesTopics(t *testing.T) {
	src := `{"exercises": [{"topics": ["APPLE","f.i$g>","honeydew      melon"]}]}`

	var srcCfg Config
	if err := json.NewDecoder(strings.NewReader(src)).Decode(&srcCfg); err != nil {
		t.Fatal(err)
	}

	dst, err := srcCfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	var dstCfg Config
	if err := json.NewDecoder(bytes.NewReader(dst)).Decode(&dstCfg); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, []string{"apple", "fig", "honeydew_melon"}, dstCfg.Exercises[0].Topics)
}

var allExercisesKeys = []string{
	"slug",
	"uuid",
}

var activeExercisesKeys = []string{
	"core",
	// not auto_approve, since it is omitempty
	"unlocked_by",
	"difficulty",
	"topics",
}

func TestMarshalActive(t *testing.T) {
	srcCfg := Config{
		Exercises: []ExerciseMetadata{
			ExerciseMetadata{
				Slug:         "active",
				Topics:       []string{"topic_one", "topic_two"},
				IsDeprecated: false,
			},
		},
	}

	dst, err := srcCfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	// contains all the keys we expect it to contain:
	for _, key := range append(allExercisesKeys, activeExercisesKeys...) {
		assert.True(t, strings.Contains(string(dst), key), "expected JSON representation to contain %q, but it didn't: %s", key, string(dst))
	}

	var dstCfg Config
	if err := json.NewDecoder(bytes.NewReader(dst)).Decode(&dstCfg); err != nil {
		t.Fatal(err)
	}

	// survived an encode -> decode:
	assert.Equal(t, "active", dstCfg.Exercises[0].Slug)
	assert.False(t, dstCfg.Exercises[0].IsDeprecated)
	assert.Equal(t, []string{"topic_one", "topic_two"}, dstCfg.Exercises[0].Topics)
}

func TestMarshalDeprecated(t *testing.T) {
	srcCfg := Config{
		Exercises: []ExerciseMetadata{
			ExerciseMetadata{
				Slug:         "deprecated",
				Topics:       []string{"topic_one", "topic_two"},
				IsDeprecated: true,
			},
		},
	}

	dst, err := srcCfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	// contains the keys we want, and not the ones we don't:
	for _, key := range append(allExercisesKeys, "deprecated") {
		assert.True(t, strings.Contains(string(dst), key), "expected JSON representation to contain %q, but it didn't: %s", key, string(dst))
	}
	for _, key := range activeExercisesKeys {
		assert.False(t, strings.Contains(string(dst), key), "expected JSON representation NOT to contain %q, but it did: %s", key, string(dst))
	}

	var dstCfg Config
	if err := json.NewDecoder(bytes.NewReader(dst)).Decode(&dstCfg); err != nil {
		t.Fatal(err)
	}

	// survived an encode -> decode:
	assert.Equal(t, "deprecated", dstCfg.Exercises[0].Slug)
	assert.True(t, dstCfg.Exercises[0].IsDeprecated)
	// Note that since topics was never marshalled, it should be nil.
	assert.Nil(t, dstCfg.Exercises[0].Topics)
}

func TestSemanticsOfMissingTopics(t *testing.T) {
	src := `
	{
		"exercises": [{
			"topics": null
		},
		{
			"topics": []
		}]
	}
	`
	var srcCfg Config
	if err := json.NewDecoder(strings.NewReader(src)).Decode(&srcCfg); err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, srcCfg.Exercises[0].Topics)
	assert.NotEqual(t, []string{}, srcCfg.Exercises[0].Topics)
	assert.NotNil(t, srcCfg.Exercises[1].Topics)
	assert.Equal(t, []string{}, srcCfg.Exercises[1].Topics)

	// Round-trip it through serialization.
	dst, err := srcCfg.ToJSON()
	if err != nil {
		t.Fatal(err)
	}

	var dstCfg Config
	if err := json.NewDecoder(bytes.NewReader(dst)).Decode(&dstCfg); err != nil {
		t.Fatal(err)
	}
	assert.Nil(t, dstCfg.Exercises[0].Topics)
	assert.NotEqual(t, []string{}, dstCfg.Exercises[0].Topics)
	assert.NotNil(t, dstCfg.Exercises[1].Topics)
	assert.Equal(t, []string{}, dstCfg.Exercises[1].Topics)
}
