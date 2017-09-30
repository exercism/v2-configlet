package track

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	// ProblemSpecificationsDir is the default name of the cloned problem-specifications repository.
	ProblemSpecificationsDir = "problem-specifications"
	filenameDescription      = "description.md"
	filenameMetadata         = "metadata.yml"
)

var (
	// ProblemSpecificationsPath is the location of the cloned problem-specifications repository.
	ProblemSpecificationsPath string
)

// ProblemSpecification contains metadata describing an exercise.
type ProblemSpecification struct {
	Slug            string
	Description     string
	Title           string `yaml:"title"`
	Source          string `yaml:"source"`
	SourceURL       string `yaml:"source_url"`
	root            string
	trackID         string
	metadataPath    string
	descriptionPath string
}

// NewProblemSpecification loads the specification from files on disk.
// It will default to a custom specification, falling back to the generic specification
// if no custom one is found.
func NewProblemSpecification(root, trackID, slug string) (*ProblemSpecification, error) {
	spec := &ProblemSpecification{
		root:    root,
		trackID: trackID,
		Slug:    slug,
	}
	spec.Title = spec.titleCasedSlug()

	if err := spec.loadMetadata(); err != nil {
		return nil, err
	}

	if err := spec.loadDescription(); err != nil {
		return nil, err
	}

	return spec, nil
}

// Name is a readable version of the slug.
func (spec *ProblemSpecification) Name() string {
	if spec.Title == "" {
		spec.Title = spec.titleCasedSlug()
	}
	return spec.Title
}

// MixedCaseName returns the name with all spaces removed.
func (spec *ProblemSpecification) MixedCaseName() string {
	return strings.Replace(spec.titleCasedSlug(), " ", "", -1)
}

// SnakeCaseName converts the slug to snake case.
func (spec *ProblemSpecification) SnakeCaseName() string {
	return strings.Replace(spec.Slug, "-", "_", -1)
}

// Credits are a markdown-formatted version of the source of the exercise.
func (spec *ProblemSpecification) Credits() string {
	if spec.SourceURL == "" {
		return spec.Source
	}
	if spec.Source == "" {
		return fmt.Sprintf("[%s](%s)", spec.SourceURL, spec.SourceURL)
	}
	return fmt.Sprintf("%s [%s](%s)", spec.Source, spec.SourceURL, spec.SourceURL)
}

func (spec *ProblemSpecification) titleCasedSlug() string {
	return strings.Title(strings.Join(strings.Split(spec.Slug, "-"), " "))
}

func (spec *ProblemSpecification) loadMetadata() error {
	metadataPath := filepath.Join(spec.customPath(), filenameMetadata)
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		metadataPath = filepath.Join(spec.sharedPath(), filenameMetadata)
	}
	spec.metadataPath = metadataPath

	b, err := ioutil.ReadFile(spec.metadataPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, &spec)
}

func (spec *ProblemSpecification) loadDescription() error {
	descriptionPath := filepath.Join(spec.customPath(), filenameDescription)
	if _, err := os.Stat(descriptionPath); os.IsNotExist(err) {
		descriptionPath = filepath.Join(spec.sharedPath(), filenameDescription)
	}
	spec.descriptionPath = descriptionPath

	b, err := ioutil.ReadFile(spec.descriptionPath)
	if err != nil {
		return err
	}
	spec.Description = string(b)

	return nil
}

func (spec *ProblemSpecification) sharedPath() string {
	if ProblemSpecificationsPath != "" {
		return filepath.Join(ProblemSpecificationsPath, "exercises", spec.Slug)
	}
	return filepath.Join(spec.root, ProblemSpecificationsDir, "exercises", spec.Slug)
}

func (spec *ProblemSpecification) customPath() string {
	return filepath.Join(spec.root, spec.trackID, "exercises", spec.Slug, ".meta")
}
