package track

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProblemSpecification(t *testing.T) {
	tests := []struct {
		description string
		slug        string
		specPath    string
		expected    ProblemSpecification
	}{
		{
			description: "shared spec if custom is missing",
			slug:        "one",
			expected: ProblemSpecification{
				Description: "This is one.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
		{
			description: "custom spec overrides shared",
			slug:        "two",
			expected: ProblemSpecification{
				Description: "This is two, customized.\n",
				Source:      "The web.",
				SourceURL:   "",
			},
		},
		{
			description: "shared spec from alternate problem-specifications location",
			slug:        "one",
			specPath:    filepath.FromSlash("../fixtures/alternate/problem-specifications"),
			expected: ProblemSpecification{
				Description: "This is the alternate one.\n",
				Source:      "The internet.",
				SourceURL:   "http://example.com",
			},
		},
	}

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
		slug  string
		name  string
		mixed string
	}{
		{
			slug:  "apple",
			name:  "Apple",
			mixed: "Apple",
		},
		{
			slug:  "apple-pie",
			name:  "Apple Pie",
			mixed: "ApplePie",
		},
		{
			slug:  "1-apple-per-day",
			name:  "1 Apple Per Day",
			mixed: "1ApplePerDay",
		},
	}

	for _, test := range tests {
		spec := ProblemSpecification{Slug: test.slug}
		assert.Equal(t, test.name, spec.Name())
		assert.Equal(t, test.mixed, spec.MixedCaseName())
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
