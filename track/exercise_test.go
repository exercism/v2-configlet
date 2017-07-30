package track

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExerciseSlug(t *testing.T) {
	path := filepath.FromSlash("../fixtures/fake-exercise")

	ex, err := NewExercise(path, PatternGroup{})
	assert.NoError(t, err)
	assert.Equal(t, "fake-exercise", ex.Slug)
}

func TestExerciseSolutionPaths(t *testing.T) {
	tests := []struct {
		PatternGroup
		path string
	}{
		{
			// It finds files in the root of the exercise directory.
			PatternGroup{SolutionPattern: "[Ee]xample"},
			"example.ext",
		},
		{
			// It finds files in a subdirectory.
			PatternGroup{SolutionPattern: "solution"},
			"subdir/solution.ext",
		},
		{
			// It only matches files, not directories.
			PatternGroup{SolutionPattern: "subdir"},
			"subdir/solution.ext",
		},
		// It finds hidden files.
		{
			PatternGroup{SolutionPattern: "secret-solution"},
			"subdir/.secret-solution.ext",
		},
		// it finds files in hidden directories
		{
			PatternGroup{SolutionPattern: "hidden.file\\.ext"},
			".hidden/file.ext",
		},
	}

	path := filepath.FromSlash("../fixtures/fake-exercise")

	for _, test := range tests {
		ex, err := NewExercise(path, test.PatternGroup)
		assert.NoError(t, err)

		assert.Equal(t, test.path, ex.SolutionPath)
	}
}
func TestExerciseTestSuitePaths(t *testing.T) {
	tests := []struct {
		PatternGroup
		path string
	}{
		{
			// It finds files in the root of the exercise directory.
			PatternGroup{TestPattern: "(?i)test"},
			"fake_test.ext",
		},
		{
			// It finds files in a subdirectory.
			PatternGroup{TestPattern: "specs"},
			"specs/file.ext",
		},
	}

	path := filepath.FromSlash("../fixtures/fake-exercise")

	for _, test := range tests {
		ex, err := NewExercise(path, test.PatternGroup)
		assert.NoError(t, err)

		assert.Equal(t, test.path, ex.TestSuitePath)
	}
}
