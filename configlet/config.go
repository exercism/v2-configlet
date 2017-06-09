package configlet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

// Config is an Exercism track configuration.
type Config struct {
	path            string
	Slug            string
	Language        string
	Active          bool
	Repository      string
	Exercises       []Exercise
	Ignored         []string
	Deprecated      []string
	Foregone        []string
	SolutionPattern string `json:"solution_pattern"`
}

// Exercise configures metadata about an implemented exercise.
// It's listed in the config in the order that the exercise will be
// delivered by the API.
type Exercise struct {
	Slug       string
	Difficulty int
	Topics     []string
}

// Load loads an Exercism track configuration.
func Load(file string) (Config, error) {
	c := NewConfig()

	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return c, fmt.Errorf("Unable to parse config: %s -- %s", file, err.Error())
	}

	return c, nil
}

// NewConfig creates a new Config with optional defaults set.
// Currently the only optional value is SolutionPattern which is used
// to work out if an exercise has a sample solution.
func NewConfig() Config {
	return Config{SolutionPattern: "[Ee]xample"}
}

// Slugs is the list of exercise identifiers for the track.
func (c Config) Slugs() []string {
	var slugs []string
	if len(c.Exercises) > 0 {
		for _, ex := range c.Exercises {
			slugs = append(slugs, ex.Slug)
		}
		return slugs
	}
	return slugs
}

// IgnoredDirs merges configured and default dirs.
// Some directories will never, ever represent an
// Exercism problem.
func (c Config) IgnoredDirs() []string {
	dirs := append(c.Ignored, "bin", "img")
	dirs = uniq(dirs)
	sort.Strings(dirs)
	return dirs
}

func uniq(items []string) []string {
	uniques := map[string]bool{}
	for _, item := range items {
		uniques[item] = true
	}

	items = []string{}
	for unique := range uniques {
		items = append(items, unique)
	}
	return items
}
