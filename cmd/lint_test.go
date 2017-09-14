package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	"github.com/exercism/configlet/track"
	"github.com/stretchr/testify/assert"
)

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
	tests := []struct {
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
	}

	for _, test := range tests {
		track := track.Track{Config: test.config}
		uuids := duplicateUUID(track)

		assert.Equal(t, test.expected, len(uuids), test.desc)
	}
}

func TestDuplicateTrackUUID(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintln(w, `{"uuids": ["ccc"]}`)
	})

	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()

	saved := UUIDValidationURL
	UUIDValidationURL = ts.URL
	defer func() { UUIDValidationURL = saved }()

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

func TestInvalidRubyRegexPatterns(t *testing.T) {
	fakeEndpoint := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintln(w, `{"patterns": ["ignore_pattern"]}`)
	})

	ts := httptest.NewServer(fakeEndpoint)
	defer ts.Close()

	saved := UUIDValidationURL
	UUIDValidationURL = ts.URL
	defer func() { UUIDValidationURL = saved }()

	expected := []string{"ignore_pattern"}
	track := track.Track{
		Config: track.Config{
			PatternGroup: track.PatternGroup{
				IgnorePattern:   "example(?!_test.go)",
				SolutionPattern: "example",
			},
		},
	}

	uuids := invalidRubyRegexPatterns(track)
	assert.Equal(t, len(expected), len(uuids))
	assert.Equal(t, expected[0], uuids[0])

}
