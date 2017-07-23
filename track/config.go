package track

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

var errInvalidConfig = errors.New("invalid config file - try jsonlint.com")

// Config is an Exercism track configuration.
type Config struct {
	Language        string
	Active          bool
	Exercises       []ExerciseMetadata
	DeprecatedSlugs []string `json:"deprecated"`
	ForegoneSlugs   []string `json:"foregone"`
	SolutionPattern string   `json:"solution_pattern"`
}

// NewConfig loads a track configuration file.
// The config has a default solution pattern if none is provided in the file.
// The solution pattern is sued to determine if an exercise has a sample solution.
func NewConfig(path string) (Config, error) {
	c := Config{
		SolutionPattern: "[Ee]xample",
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
