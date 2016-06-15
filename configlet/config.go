package configlet

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
)

// Config is an Exercism track configuration.
type Config struct {
	path       string
	Slug       string
	Language   string
	Active     bool
	Repository string
	Problems   []string
	Ignored    []string
	Deprecated []string
	Foregone   []string
}

// Load loads an Exercism track configuration.
func Load(file string) (Config, error) {
	c := Config{}

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

// IgnoredDirs merges configured and default dirs.
// Some directories will never, ever represent an
// Exercism problem.
func (c Config) IgnoredDirs() []string {
	dirs := append(c.Ignored, "bin")
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
