package track

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProblemSpecification(t *testing.T) {
	tests := []struct {
		desc      string
		slug      string
		trackPath string
		specPath  string
		expected  ProblemSpecification
	}{
		{
			desc:      "shared spec if custom is missing",
			slug:      "one",
			trackPath: filepath.FromSlash("../fixtures/numbers"),
			expected: ProblemSpecification{
				Description: "This is one.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
		{
			desc:      "custom spec overrides shared",
			slug:      "two",
			trackPath: filepath.FromSlash("../fixtures/numbers"),
			expected: ProblemSpecification{
				Description: "This is two, customized.\n",
				Source:      "The web.",
				SourceURL:   "",
			},
		},
		{
			desc:      "shared spec from alternate problem-specifications location",
			slug:      "one",
			trackPath: filepath.FromSlash("../fixtures/numbers"),
			specPath:  filepath.FromSlash("../fixtures/alternate/problem-specifications"),
			expected: ProblemSpecification{
				Description: "This is the alternate one.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
		{
			desc:      "custom spec metadata with shared description",
			slug:      "metadata-example",
			trackPath: filepath.FromSlash("../fixtures/granular-metadata-overrides"),
			specPath:  filepath.FromSlash("../fixtures/granular-metadata-overrides/problem-specifications"),
			expected: ProblemSpecification{
				Description: "This is a shared description.\n",
				Source:      "The web.",
				SourceURL:   "",
			},
		},
		{
			desc:      "shared spec metadata with custom description",
			slug:      "description-example",
			trackPath: filepath.FromSlash("../fixtures/granular-metadata-overrides"),
			specPath:  filepath.FromSlash("../fixtures/granular-metadata-overrides/problem-specifications"),
			expected: ProblemSpecification{
				Description: "This is a custom description.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
	}
	originalSpecPath := ProblemSpecificationsPath
	defer func() { ProblemSpecificationsPath = originalSpecPath }()

	for _, test := range tests {
		ProblemSpecificationsPath = test.specPath
		root, trackID := filepath.Dir(test.trackPath), filepath.Base(test.trackPath)
		spec, err := NewProblemSpecification(root, trackID, test.slug)
		assert.NoError(t, err)

		assert.Equal(t, test.expected.Source, spec.Source)
		assert.Equal(t, test.expected.SourceURL, spec.SourceURL)
		assert.Equal(t, test.expected.Description, spec.Description)
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
	root := filepath.FromSlash("../fixtures")
	originalSpecPath := ProblemSpecificationsPath
	ProblemSpecificationsPath = filepath.Join(root, "titled-problem-specifications")
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

	for _, test := range tests {
		spec, err := NewProblemSpecification(root, "titled-problem-specifications", test.slug)
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
