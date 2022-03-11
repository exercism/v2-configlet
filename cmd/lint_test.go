package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"testing"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
	"github.com/stretchr/testify/assert"
)

var (
	// The JSON needs to have nulls, not empty strings.
	// We cannot take the address of strings.
	apple   = "apple"
	banana  = "banana"
	unknown = "unknown"
)

func TestLintTrack(t *testing.T) {
	originalNoHTTP := noHTTP
	noHTTP = true
	defer func() {
		noHTTP = originalNoHTTP
	}()

	originalOut := ui.Out
	originalErrOut := ui.ErrOut
	ui.Out = ioutil.Discard
	ui.ErrOut = ioutil.Discard
	defer func() {
		ui.Out = originalOut
		ui.ErrOut = originalErrOut
	}()

	lintTests := []struct {
		desc     string
		path     string
		expected bool
	}{
		{
			desc:     "should fail when a track contains one or more lint failures.",
			path:     "../fixtures/numbers",
			expected: true,
		},
		{
			desc:     "should fail when a track contains malformed configuration data.",
			path:     "../fixtures/broken-maintainers",
			expected: true,
		},
		{
			desc:     "should fail when a track is missing a README.",
			path:     "../fixtures/missing-readme",
			expected: true,
		},
		{
			desc:     "should not fail when given a track with all of its bits in place.",
			path:     "../fixtures/lint/valid-track",
			expected: false,
		},
	}

	for _, tt := range lintTests {
		failed := lintTrack(filepath.FromSlash(tt.path))
		assert.Equal(t, tt.expected, failed, tt.desc)
	}
}

func TestMissingImplementations(t *testing.T) {
	track := track.Track{
		Config: track.Config{
			Exercises: []track.ExerciseMetadata{
				{Slug: "apple"},
				{Slug: "banana"},
				{Slug: "cherry"},
			},
		},
		Exercises: []track.Exercise{
			{Slug: "apple"},
		},
	}

	slugs := missingImplementations(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected missing implementations in 2 exercises, missing in %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "banana", slugs[0])
	assert.Equal(t, "cherry", slugs[1])
}

func TestMissingMetadata(t *testing.T) {
	track := track.Track{
		Config: track.Config{
			Exercises: []track.ExerciseMetadata{
				{Slug: "apple"},
			},
		},
		Exercises: []track.Exercise{
			{Slug: "apple"},
			{Slug: "banana"},
			{Slug: "cherry"},
		},
	}

	slugs := missingMetadata(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected missing metadata in 2 exercises, missing in %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "banana", slugs[0])
	assert.Equal(t, "cherry", slugs[1])
}

func TestMissingReadme(t *testing.T) {
	track := track.Track{
		Exercises: []track.Exercise{
			{Slug: "apple"},
			{Slug: "banana", ReadmePath: "README.md"},
			{Slug: "cherry"},
		},
	}

	slugs := missingReadme(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected missing READMEs in 2 exercises, missing in %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "apple", slugs[0])
	assert.Equal(t, "cherry", slugs[1])
}

func TestMissingSolution(t *testing.T) {
	track := track.Track{
		Exercises: []track.Exercise{
			{Slug: "apple"},
			{Slug: "banana", SolutionPath: "b.txt"},
			{Slug: "cherry"},
		},
	}

	slugs := missingSolution(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected missing solutions in 2 exercises, missing in %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "apple", slugs[0])
	assert.Equal(t, "cherry", slugs[1])
}

func TestMissingTestSuite(t *testing.T) {
	track := track.Track{
		Exercises: []track.Exercise{
			{Slug: "apple"},
			{Slug: "banana", TestSuitePath: "b_test.ext"},
			{Slug: "cherry"},
		},
	}

	slugs := missingTestSuite(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected missing test in 2 exercises, missing in %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "apple", slugs[0])
	assert.Equal(t, "cherry", slugs[1])
}

func TestForegoneViolations(t *testing.T) {
	track := track.Track{
		Config: track.Config{
			ForegoneSlugs: []string{"banana", "cherry"},
		},
		Exercises: []track.Exercise{
			{Slug: "apple"},
			{Slug: "banana"},
			{Slug: "cherry"},
		},
	}

	slugs := foregoneViolations(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected foregone violations in 2 exercises, violations in %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "banana", slugs[0])
	assert.Equal(t, "cherry", slugs[1])
}

func TestDuplicateSlugs(t *testing.T) {
	track := track.Track{
		Config: track.Config{
			Exercises: []track.ExerciseMetadata{
				{Slug: "apple"},
				{Slug: "banana"},
				{Slug: "cherry"},
			},
			DeprecatedSlugs: []string{"apple"},
			ForegoneSlugs:   []string{"banana"},
		},
	}

	slugs := duplicateSlugs(track)

	if len(slugs) != 2 {
		t.Fatalf("Expected 2 duplicate slugs, found %d", len(slugs))
	}

	sort.Strings(slugs)

	assert.Equal(t, "apple", slugs[0])
	assert.Equal(t, "banana", slugs[1])
}

func TestDuplicateUUID(t *testing.T) {
	uuidTests := []struct {
		desc     string
		expected int
		config   track.Config
	}{
		{
			desc:     "should not complain about a conflicting UUID for exercises with missing UUIDs.",
			expected: 0,
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{Slug: "apple", UUID: ""},
					{Slug: "banana", UUID: ""},
				},
			},
		},
		{
			desc:     "should complain that multiple exercises have a conflicting UUID.",
			expected: 1,
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{Slug: "cherry", UUID: "ccc"},
					{Slug: "diakon", UUID: "abc"},
					{Slug: "eggplant", UUID: "ccc"},
				},
			},
		},
		{
			desc:     "should complain that multiple exercises have a conflicting UUID (while trimming spaces).",
			expected: 1,
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{Slug: "cherry", UUID: " ccc"},
					{Slug: "diakon", UUID: "abc"},
					{Slug: "eggplant", UUID: "ccc "},
				},
			},
		},
	}

	for _, tt := range uuidTests {
		track := track.Track{Config: tt.config}
		uuids := duplicateUUID(track)

		assert.Equal(t, tt.expected, len(uuids), tt.desc)
	}
}

func TestDuplicateTrackUUID(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, `{"uuids": ["ccc"]}`)
	})

	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()

	originalUUIDValidationURL := UUIDValidationURL
	UUIDValidationURL = ts.URL
	defer func() { UUIDValidationURL = originalUUIDValidationURL }()

	expected := []string{"ccc"}
	track := track.Track{
		Config: track.Config{
			Exercises: []track.ExerciseMetadata{
				{Slug: "apple", UUID: "abc"},
				{Slug: "banana", UUID: expected[0]},
			},
		},
	}

	uuids := duplicateTrackUUID(track)
	assert.Equal(t, len(expected), len(uuids))
	assert.Equal(t, expected[0], uuids[0])

}

func TestLockedCoreViolation(t *testing.T) {
	trackTests := []struct {
		desc               string
		expectedViolations int
		expectedSlug       string
		config             track.Config
	}{
		{
			desc:               "should have one locked core violation",
			expectedViolations: 1,
			expectedSlug:       "apple",
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{
						Slug:       "apple",
						UnlockedBy: &banana,
						IsCore:     true,
					},
					{
						Slug:       "banana",
						UnlockedBy: nil,
						IsCore:     false,
					},
				},
			},
		},
		{
			desc:               "should have no locked core violations",
			expectedViolations: 0,
			expectedSlug:       "",
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{
						Slug:       "apple",
						UnlockedBy: nil,
						IsCore:     true,
					},
					{
						Slug:       "banana",
						UnlockedBy: &apple,
						IsCore:     false,
					},
					{
						Slug:       "cherry",
						UnlockedBy: nil,
						IsCore:     false,
					},
				},
			},
		},
	}

	for _, tt := range trackTests {
		track := track.Track{Config: tt.config}
		slugs := lockedCoreViolation(track)

		assert.Equal(t, tt.expectedViolations, len(slugs), tt.desc)

		if len(slugs) > 0 {
			assert.Equal(t, tt.expectedSlug, slugs[0], tt.desc)
		}
	}
}

func TestUnlockedByViolations(t *testing.T) {
	trackTests := []struct {
		desc               string
		expectedViolations int
		expectedSlug       string
		config             track.Config
	}{
		{
			desc:               "should have one unlocked by violation from non-core exercise",
			expectedViolations: 1,
			expectedSlug:       "banana",
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{
						Slug:       "apple",
						UnlockedBy: nil,
						IsCore:     false,
					},
					{
						Slug:       "banana",
						UnlockedBy: &apple,
						IsCore:     false,
					},
				},
			},
		},
		{
			desc:               "should have one unlocked by violation from non-existent exercise",
			expectedViolations: 1,
			expectedSlug:       "cherry",
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{
						Slug:       "apple",
						UnlockedBy: nil,
						IsCore:     true,
					},
					{
						Slug:       "banana",
						UnlockedBy: &apple,
						IsCore:     false,
					},
					{
						Slug:       "cherry",
						UnlockedBy: &unknown,
						IsCore:     false,
					},
				},
			},
		},
		{
			desc:               "should have no unlocked by violations",
			expectedViolations: 0,
			expectedSlug:       "",
			config: track.Config{
				Exercises: []track.ExerciseMetadata{
					{
						Slug:       "apple",
						UnlockedBy: nil,
						IsCore:     true,
					},
					{
						Slug:       "banana",
						UnlockedBy: &apple,
						IsCore:     false,
					},
				},
			},
		},
	}

	for _, tt := range trackTests {
		track := track.Track{Config: tt.config}
		slugs := unlockedByValidExercise(track)

		assert.Equal(t, tt.expectedViolations, len(slugs), tt.desc)

		if len(slugs) > 0 {
			assert.Equal(t, tt.expectedSlug, slugs[0], tt.desc)
		}
	}
}
