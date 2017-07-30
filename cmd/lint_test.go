package cmd

import (
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
