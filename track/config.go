package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// PatternGroup holds matching patterns defined in an Exercism track config
type PatternGroup struct {
	IgnorePattern   string `json:"ignore_pattern,omitempty"`
	SolutionPattern string `json:"solution_pattern,omitempty"`
	TestPattern     string `json:"test_pattern,omitempty"`
}

// ConceptKey defines the fields associated with the full array of concepts
type ConceptKey struct {
	Slug  string `json:"slug"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Blurb string `json:"blurb"`
}

// PracticeMetadata contains metadata about an implemented practice exercise.
// The order below determines the order delivered by the API.
type PracticeMetadata struct {
	Slug         string   `json:"slug"`
	UUID         string   `json:"uuid"`
	IsCore       bool     `json:"core"`
	AutoApprove  bool     `json:"auto_approve,omitempty"`
	UnlockedBy   *string  `json:"unlocked_by"`
	Difficulty   int      `json:"difficulty"`
	Topics       []string `json:"topics,omitempty"`
	IsDeprecated bool     `json:"deprecated,omitempty"`
}

// ConceptMetadata contains metadata about an implemented concept exercise.
// The order below determines the order delivered by the API.
type ConceptMetadata struct {
	Slug         string   `json:"slug"`
	UUID         string   `json:"uuid"`
	IsCore       bool     `json:"core"`
	AutoApprove  bool     `json:"auto_approve,omitempty"`
	UnlockedBy   *string  `json:"unlocked_by"`
	Difficulty   int      `json:"difficulty"`
	Topics       []string `json:"topics,omitempty"`
	IsDeprecated bool     `json:"deprecated,omitempty"`
}

// Config is an Exercism track configuration.
type Config struct {
	TrackID      string `json:"track_id,omitempty"`
	Version      int    `json:"version"`
	OnlineEditor struct {
		IndentStyle string `json:"indent_style"`
		IndentSize  int    `json:"indent_size"`
	} `json:"online_editor"`
	Language       string `json:"language"`
	Active         bool   `json:"active"`
	Blurb          string `json:"blurb"`
	Gitter         string `json:"gitter,omitempty"`
	ChecklistIssue int    `json:"checklist_issue,omitempty"`
	PatternGroup
	ForegoneSlugs []string     `json:"foregone,omitempty"`
	Concepts      []ConceptKey `json:"concepts"`
	Exercises     struct {
		ConceptExercises  []ConceptMetadata  `json:"concept"`
		PracticeExercises []PracticeMetadata `json:"exercises"`
	} `json:"exercises"`
	DeprecatedSlugs []string `json:"deprecated,omitempty"`
}

// NewConfig loads a track configuration file.
// The config has a default solution and test pattern if not provided in the file.
// The solution pattern is used to determine if an exercise has a sample solution.
// The test pattern is used to determine if an exercise has a test suite.
// The ignore pattern is used to exclude files from the 'exercism fetch' command.
func NewConfig(path string) (Config, error) {
	c := Config{
		PatternGroup: PatternGroup{
			SolutionPattern: "[Ee]xample",
			TestPattern:     "(?i)test",
			IgnorePattern:   "[Ee]xample",
		},
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return c, fmt.Errorf("invalid config %s -- %s", path, err.Error())
	}
	return c, nil
}

// LoadFromFile loads a config from file given the path to the file.
func (cfg *Config) LoadFromFile(path string) error {
	file, err := os.Open(filepath.FromSlash(path))
	if err != nil {
		return err
	}

	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return err
	}
	return nil
}

// ToJSON marshals the Config to normalized JSON.
func (cfg *Config) ToJSON() ([]byte, error) {
	return json.MarshalIndent(&cfg, "", "  ")
}
