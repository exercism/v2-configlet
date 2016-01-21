package configlet

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// Track represents a set of Exercism problems.
// Typically these will all be in the same language, defined
// in github.com/exercism/x<LANGUAGE>.
type Track struct {
	path string
}

// NewTrack finds a Track at path.
// It will look for a config.json in the root of that path.
// This file will list problems that correspond to
// directories which contain a test suite and supporting
// files, along with an example solution.
func NewTrack(path string) Track {
	return Track{path: path}
}

// Config loads a track's configuration.
func (t Track) Config() (Config, error) {
	c, err := Load(t.configFile())
	if err != nil {
		return c, err
	}
	return c, nil
}

// HasValidConfig lints the JSON file.
func (t Track) HasValidConfig() bool {
	_, err := t.Config()
	return err == nil
}

// Problems lists all the configured problems.
func (t Track) Problems() (map[string]struct{}, error) {
	problems := make(map[string]struct{})

	c, err := t.Config()
	if err != nil {
		return problems, err
	}

	for _, problem := range c.Problems {
		problems[problem] = struct{}{}
	}

	return problems, nil
}

// Slugs is a list of all problems mentioned in the config.
func (t Track) Slugs() (map[string]struct{}, error) {
	slugs := make(map[string]struct{})

	c, err := t.Config()
	if err != nil {
		return slugs, err
	}

	for _, slug := range c.Problems {
		slugs[slug] = struct{}{}
	}

	for _, slug := range c.IgnoredDirs() {
		slugs[slug] = struct{}{}
	}

	for _, slug := range c.Deprecated {
		slugs[slug] = struct{}{}
	}

	for _, slug := range c.Foregone {
		slugs[slug] = struct{}{}
	}
	return slugs, nil
}

// Dirs is a list of all the relevant directories.
func (t Track) Dirs() (map[string]struct{}, error) {
	dirs := make(map[string]struct{})

	infos, err := ioutil.ReadDir(t.path)
	if err != nil {
		return dirs, err
	}

	for _, info := range infos {
		if info.IsDir() && info.Name() != "exercises" {
			dirs[info.Name()] = struct{}{}
		}
	}

	path := filepath.Join(t.path, "exercises")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return dirs, nil
		}
		return dirs, err
	}

	infos, err = ioutil.ReadDir(filepath.Join(t.path, "exercises"))
	if err != nil {
		return dirs, err
	}

	for _, info := range infos {
		if info.IsDir() && info.Name() != "exercises" {
			dirs[info.Name()] = struct{}{}
		}
	}

	return dirs, nil
}

// MissingProblems identify problems lacking an implementation.
// This will complain if the problem slug is listed in the configuration,
// but there is no corresponding directory for it.
func (t Track) MissingProblems() ([]string, error) {
	dirs, err := t.Dirs()
	if err != nil {
		return []string{}, err
	}

	problems, err := t.Problems()
	if err != nil {
		return []string{}, err
	}

	omissions := make([]string, 0, len(problems))

	for problem := range problems {
		_, present := dirs[problem]
		if !present {
			omissions = append(omissions, problem)
		}
	}
	return omissions, nil
}

// UnconfiguredProblems identifies unlisted implementations.
// This will complain if a directory exists, but is not mentioned
// anywhere in the config file.
func (t Track) UnconfiguredProblems() ([]string, error) {
	dirs, err := t.Dirs()
	if err != nil {
		return []string{}, err
	}

	slugs, err := t.Slugs()
	if err != nil {
		return []string{}, err
	}

	omissions := make([]string, 0, len(slugs))

	for dir := range dirs {
		_, present := slugs[dir]
		if !present {
			omissions = append(omissions, dir)
		}
	}
	return omissions, nil
}

// ProblemsLackingExample identifies implementations without a solution.
// This will often be triggered because the implementation's sample solution
// is not named something with example. This is particularly critical since
// any file that is in a path not named /[Ee]xample/ will be served by the API,
// showing the user a possible solution before they have solved the problem
// themselves.
func (t Track) ProblemsLackingExample() ([]string, error) {
	problems := []string{}

	c, err := t.Config()
	if err != nil {
		return problems, err
	}

	for _, problem := range c.Problems {
		filename := fmt.Sprintf("%s/%s", t.path, problem)
		if _, err := os.Stat(filename); err == nil {
			files, err := findAllFiles(fmt.Sprintf("%s/%s", t.path, problem))
			if err != nil {
				return problems, err
			}
			found, err := hasExampleFile(files)
			if !found {
				problems = append(problems, problem)
			}
		}
	}

	return problems, nil
}

// ForegoneViolations indentifies implementations that should not be included.
// This could be because the problem is too trivial, ridiculously non-trivial,
// or simply uninteresting.
func (t Track) ForegoneViolations() ([]string, error) {
	problems := []string{}

	c, err := t.Config()
	if err != nil {
		return problems, err
	}

	dirs, err := t.Dirs()
	if err != nil {
		return problems, err
	}

	violations := make([]string, 0, len(dirs))

	for _, problem := range c.Foregone {
		_, present := dirs[problem]
		if present {
			violations = append(violations, problem)
		}
	}
	return violations, nil
}

// DuplicateSlugs detects slugs in multiple config categories.
// If a problem is deprecated, it means that we have the files for it,
// we're just not serving it in the default response.
// If a directory is ignored, it means that it's not a problem.
// If a slug is foregone, it means that we've chosen not to implement it,
// and it should not have a directory.
func (t Track) DuplicateSlugs() ([]string, error) {
	counts := make(map[string]int)

	c, err := t.Config()
	if err != nil {
		return []string{}, err
	}

	for _, slug := range c.Problems {
		counts[slug] = counts[slug] + 1
	}

	for _, slug := range c.IgnoredDirs() {
		counts[slug] = counts[slug] + 1
	}

	for _, slug := range c.Deprecated {
		counts[slug] = counts[slug] + 1
	}

	for _, slug := range c.Foregone {
		counts[slug] = counts[slug] + 1
	}

	dupes := make([]string, 0, len(counts))
	for slug, count := range counts {
		if count > 1 {
			dupes = append(dupes, slug)
		}
	}
	sort.Strings(dupes)

	return dupes, nil
}

func (t Track) configFile() string {
	return fmt.Sprintf("%s/config.json", t.path)
}

func hasExampleFile(files []string) (bool, error) {
	r, err := regexp.Compile(`[Ee]xample`)
	if err != nil {
		return false, err
	}
	for _, file := range files {
		matches := r.Find([]byte(file))
		if len(matches) > 0 {
			return true, nil
		}
	}
	return false, nil
}

func findAllFiles(path string) ([]string, error) {
	files := []string{}

	infos, err := ioutil.ReadDir(path)
	if err != nil {
		return files, err
	}

	for _, info := range infos {
		subPath := fmt.Sprintf("%s/%s", path, info.Name())
		if info.IsDir() {
			subFiles, err := findAllFiles(subPath)
			if err != nil {
				return files, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, subPath)
		}
	}
	return files, nil
}
