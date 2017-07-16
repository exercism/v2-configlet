package track

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

const (
	filenameDescription = "description.md"
	filenameMetadata    = "metadata.yml"
)

var (
	// ProblemSpecificationsPath is the location of the cloned problem-specifications repository.
	ProblemSpecificationsPath string
)

type ProblemSpecification struct {
	Slug            string
	Description     string
	Source          string `yaml:"source"`
	SourceURL       string `yaml:"source_url"`
	root            string
	trackID         string
	metadataPath    string
	descriptionPath string
	specPath        string
}

func NewProblemSpecification(root, trackID, slug string) (*ProblemSpecification, error) {
	spec := &ProblemSpecification{
		root:    root,
		trackID: trackID,
		Slug:    slug,
	}
	err := spec.load(spec.customPath())
	if err == nil {
		return spec, nil
	}
	err = spec.load(spec.sharedPath())
	if err == nil {
		return spec, nil
	}
	return nil, err
}

// Name is a readable version of the slug.
func (spec *ProblemSpecification) Name() string {
	return strings.Title(strings.Join(strings.Split(spec.Slug, "-"), " "))
}

func (spec *ProblemSpecification) Credits() string {
	if spec.SourceURL == "" {
		return spec.Source
	}
	if spec.Source == "" {
		return fmt.Sprintf("[%s](%s)", spec.SourceURL, spec.SourceURL)
	}
	return fmt.Sprintf("%s [%s](%s)", spec.Source, spec.SourceURL, spec.SourceURL)
}

func (spec *ProblemSpecification) load(path string) error {
	spec.descriptionPath = filepath.Join(path, filenameDescription)
	spec.metadataPath = filepath.Join(path, filenameMetadata)

	b, err := ioutil.ReadFile(spec.descriptionPath)
	if err != nil {
		return err
	}
	spec.Description = string(b)

	b, err = ioutil.ReadFile(spec.metadataPath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, &spec)
}

func (spec *ProblemSpecification) sharedPath() string {
	if ProblemSpecificationsPath != "" {
		return filepath.Join(ProblemSpecificationsPath, "exercises", spec.Slug)
	}
	return filepath.Join(spec.root, "problem-specifications", "exercises", spec.Slug)
}

func (spec *ProblemSpecification) customPath() string {
	return filepath.Join(spec.root, spec.trackID, "exercises", spec.Slug, ".meta")
}
