package track

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// PatternGroup holds matching patterns defined in an Exercism track configuration.
type PatternGroup struct {
	SolutionPattern string `json:"solution_pattern"`
	TestPattern     string `json:"test_pattern"`
	IgnorePattern   string `json:"ignore_pattern"`
}

// Config is an Exercism track configuration.
type Config struct {
	Language        string
	Active          bool
	Exercises       []ExerciseMetadata
	DeprecatedSlugs []string `json:"deprecated"`
	ForegoneSlugs   []string `json:"foregone"`
	PatternGroup
}

// NewConfig loads a track configuration file.
// The config has a default solution and test pattern if not provided in the file.
// The solution pattern is used to determine if an exercise has a sample solution.
// The test pattern is used to determine if an exercise has a test suite.
func NewConfig(path string) (Config, error) {
	c := Config{
		PatternGroup: PatternGroup{
			IgnorePattern:   "(?i)[Ee]xample",
			SolutionPattern: "[Ee]xample",
			TestPattern:     "(?i)test",
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

func (c Config) Patterns() []string {
	return []string{
		c.IgnorePattern,
		c.SolutionPattern,
		c.TestPattern,
	}
}
