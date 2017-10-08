package track

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProblemSpecification(t *testing.T) {
	tests := []struct {
		desc     string
		slug     string
		specPath string
		expected ProblemSpecification
	}{
		{
			desc: "shared spec if custom is missing",
			slug: "one",
			expected: ProblemSpecification{
				Description: "This is one.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
		{
			desc: "custom spec overrides shared",
			slug: "two",
			expected: ProblemSpecification{
				Description: "This is two, customized.\n",
				Source:      "The web.",
				SourceURL:   "",
			},
		},
		{
			desc:     "shared spec from alternate problem-specifications location",
			slug:     "one",
			specPath: filepath.FromSlash("../fixtures/alternate/problem-specifications"),
			expected: ProblemSpecification{
				Description: "This is the alternate one.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
	}

	originalSpecPath := ProblemSpecificationsPath
	defer func() { ProblemSpecificationsPath = originalSpecPath }()

	for _, test := range tests {
		ProblemSpecificationsPath = test.specPath
		root := filepath.FromSlash("../fixtures")
		spec, err := NewProblemSpecification(root, "numbers", test.slug)
		assert.NoError(t, err)

		assert.Equal(t, test.expected.Description, spec.Description)
		assert.Equal(t, test.expected.Source, spec.Source)
		assert.Equal(t, test.expected.SourceURL, spec.SourceURL)
	}
}

func TestMissingProblemSpecification(t *testing.T) {
	root := filepath.FromSlash("../fixtures")
	_, err := NewProblemSpecification(root, "numbers", "three")
	assert.Error(t, err)
}

func TestProblemSpecificationName(t *testing.T) {
	tests := []struct {
		desc     string
		slug     string
		title    string
		expected string
	}{
		{
			desc:     "simple slug as name",
			slug:     "apple",
			expected: "Apple",
		},
		{
			desc:     "multi-word slug as name",
			slug:     "1-apple-per-day",
			expected: "1 Apple Per Day",
		},
		{
			desc:     "title overrides slug as name",
			slug:     "rna-transcription",
			title:    "RNA Transcription",
			expected: "RNA Transcription",
		},
	}

	for _, test := range tests {
		spec := ProblemSpecification{Slug: test.slug, Title: test.title}
		assert.Equal(t, test.expected, spec.Name(), test.desc)
	}
}

func TestProblemSpecificationMixedCaseName(t *testing.T) {
	tests := []struct {
		slug     string
		expected string
	}{
		{
			slug:     "apple",
			expected: "Apple",
		},
		{
			slug:     "1-apple-per-day",
			expected: "1ApplePerDay",
		},
		{
			slug:     "rna-transcription",
			expected: "RnaTranscription",
		},
	}

	for _, test := range tests {
		spec := ProblemSpecification{Slug: test.slug}
		assert.Equal(t, test.expected, spec.MixedCaseName())
	}
}

func TestProblemSpecificationSnakeCaseName(t *testing.T) {
	tests := []struct {
		slug     string
		expected string
	}{
		{
			slug:     "apple",
			expected: "apple",
		},
		{
			slug:     "1-apple-per-day",
			expected: "1_apple_per_day",
		},
		{
			slug:     "rna-transcription",
			expected: "rna_transcription",
		},
	}

	for _, test := range tests {
		spec := ProblemSpecification{Slug: test.slug}
		assert.Equal(t, test.expected, spec.SnakeCaseName())
	}
}

func TestProblemSpecificationTitle(t *testing.T) {
	originalSpecPath := ProblemSpecificationsPath
	ProblemSpecificationsPath = ""
	defer func() { ProblemSpecificationsPath = originalSpecPath }()

	tests := []struct {
		desc     string
		slug     string
		expected string
	}{
		{
			desc:     "title is inferred from slug if not explicitly provided",
			slug:     "slug-as-title",
			expected: "Slug As Title",
		},
		{
			desc:     "explicit title takes precedence over slug",
			slug:     "explicit-title",
			expected: "(Explicit) Title",
		},
	}

	root := filepath.FromSlash("../fixtures")
	for _, test := range tests {
		spec, err := NewProblemSpecification(root, "titled-exercises", test.slug)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, spec.Title)
	}
}

func TestProblemSpecificationCredits(t *testing.T) {
	tests := []struct {
		spec    ProblemSpecification
		credits string
	}{
		{
			spec: ProblemSpecification{
				Source:    "Apple.",
				SourceURL: "http://apple.com",
			},
			credits: "Apple. [http://apple.com](http://apple.com)",
		},
		{
			spec: ProblemSpecification{
				Source:    "banana",
				SourceURL: "",
			},
			credits: "banana",
		},
		{
			spec: ProblemSpecification{
				Source:    "",
				SourceURL: "http://cherry.com",
			},
			credits: "[http://cherry.com](http://cherry.com)",
		},
		{
			spec: ProblemSpecification{
				Source:    "",
				SourceURL: "",
			},
			credits: "",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.credits, test.spec.Credits())
	}
}
