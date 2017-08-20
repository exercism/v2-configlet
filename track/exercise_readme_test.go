package track

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewExerciseReadme(t *testing.T) {
	root := filepath.FromSlash("../fixtures")

	readme, err := NewExerciseReadme(root, "numbers", "one")
	assert.NoError(t, err)
	assert.Equal(t, "This is one.\n", readme.Spec.Description)
	assert.Equal(t, "", readme.Hints)
	assert.Equal(t, "Track insert.\n", readme.TrackInsert)
	assert.Equal(t, "The {{ .Spec.Name }} exercise (from shared template).\n", readme.template)

	readme, err = NewExerciseReadme(root, "numbers", "two")
	assert.NoError(t, err)
	assert.Equal(t, "This is two, customized.\n", readme.Spec.Description)
	assert.Equal(t, "Hinting about two.\n", readme.Hints)
	assert.Equal(t, "Track insert.\n", readme.TrackInsert)
	assert.Equal(t, "{{ .Spec.Name }} has its own template with description:\n\n{{ .Spec.Description }}\n", readme.template)
}

func TestGenerateExerciseReadme(t *testing.T) {
	readme := ExerciseReadme{
		Spec: &ProblemSpecification{
			Slug:        "hello-kitty",
			Description: "The description.\n",
		},
		template: "# {{ .Spec.Name }}\n\n{{ .Spec.Description -}}",
	}
	expected := "# Hello Kitty\n\nThe description.\n"

	s, err := readme.Generate()
	assert.NoError(t, err)
	assert.Equal(t, expected, s)
}

func TestExerciseReadmeTrackInsertDeprecation(t *testing.T) {
	root := filepath.FromSlash("../fixtures/deprecated")

	tests := []struct {
		trackID  string
		expected string
	}{
		{"inserts-both", "real insert\n"},
		{"inserts-old", "deprecated insert\n"},
	}

	ProblemSpecificationsPath = filepath.FromSlash("../fixtures/problem-specifications")
	for _, test := range tests {
		readme, err := NewExerciseReadme(root, test.trackID, "fake")
		assert.NoError(t, err)
		assert.Equal(t, test.expected, readme.TrackInsert)
	}
}

func TestExerciseReadmeHintsDeprecation(t *testing.T) {
	root := filepath.FromSlash("../fixtures/deprecated")

	tests := []struct {
		trackID  string
		expected string
	}{
		{"hints-both", "real hints\n"},
		{"hints-old", "deprecated hints\n"},
	}

	ProblemSpecificationsPath = filepath.FromSlash("../fixtures/problem-specifications")
	for _, test := range tests {
		readme, err := NewExerciseReadme(root, test.trackID, "fake")
		assert.NoError(t, err)
		assert.Equal(t, test.expected, readme.Hints)
	}
}
