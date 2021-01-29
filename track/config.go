package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	rgxFunkyChars = regexp.MustCompile(`[^a-z\s-_]+`)
	rgxSpaces     = regexp.MustCompile(`[\s-]+`)
)

// PatternGroup holds matching patterns defined in an Exercism track configuration.
type PatternGroup struct {
	IgnorePattern   string `json:"ignore_pattern,omitempty"`
	SolutionPattern string `json:"solution_pattern,omitempty"`
	TestPattern     string `json:"test_pattern,omitempty"`
}

// ExerciseMetadata contains metadata about an implemented exercise.
// It's listed in the config in the order that the exercise will be
// delivered by the API.
type ExerciseMetadata struct {
	Slug         string   `json:"slug"`
	UUID         string   `json:"uuid"`
	IsCore       bool     `json:"core"`
	AutoApprove  bool     `json:"auto_approve,omitempty"`
	UnlockedBy   *string  `json:"unlocked_by"`
	Difficulty   int      `json:"difficulty"`
	Topics       []string `json:"topics"`
	IsDeprecated bool     `json:"deprecated,omitempty"`
}

// Config is an Exercism track configuration.
type Config struct {
	TrackID        string `json:"track_id,omitempty"`
	Language       string `json:"language"`
	Active         bool   `json:"active"`
	Blurb          string `json:"blurb"`
	Gitter         string `json:"gitter,omitempty"`
	ChecklistIssue int    `json:"checklist_issue,omitempty"`
	PatternGroup
	ForegoneSlugs   []string           `json:"foregone,omitempty"`
	Exercises       []ExerciseMetadata `json:"exercises"`
	DeprecatedSlugs []string           `json:"deprecated,omitempty"`
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
	for _, exercise := range cfg.Exercises {
		for i, t := range exercise.Topics {
			exercise.Topics[i] = normalizeTopic(t)
		}
		sort.Strings(exercise.Topics)
	}
	return json.MarshalIndent(&cfg, "", "  ")
}

// MarshalJSON marshals a exercise metadata to JSON,
// only marshalling certain fields if the exercise is deprecated.
func (e *ExerciseMetadata) MarshalJSON() ([]byte, error) {
	if e.IsDeprecated {
		// Only marshal Slug, UUID, Deprecated.
		return json.Marshal(&struct {
			Slug         string `json:"slug"`
			UUID         string `json:"uuid"`
			IsDeprecated bool   `json:"deprecated"`
		}{
			Slug:         e.Slug,
			UUID:         e.UUID,
			IsDeprecated: true,
		})
	} else {
		// Use the default marshalling.
		// We can't embed ExerciseMetadata into an anonymous struct,
		// since that will cause infinite recursion on this MarshalJSON,
		// But we can embed a new typedef of it,
		// since the typedef does not have this MarshalJSON function.
		// Technique discovered from http://choly.ca/post/go-json-marshalling/
		type ExerciseMetadataJ ExerciseMetadata
		return json.Marshal(&struct {
			*ExerciseMetadataJ
		}{
			ExerciseMetadataJ: (*ExerciseMetadataJ)(e),
		})
	}
}

func normalizeTopic(t string) string {
	s := strings.ToLower(t)
	s = rgxFunkyChars.ReplaceAllString(s, "")
	s = rgxSpaces.ReplaceAllString(s, "_")
	return s
}
