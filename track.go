package main

import (
	"fmt"
	"io/ioutil"
)

type Track struct {
	path string
}

func NewTrack(path string) Track {
	return Track{path: path}
}

func (t Track) configFile() string {
	return fmt.Sprintf("%s/config.json", t.path)
}

func (t Track) Config() (Config, error) {
	c, err := Load(t.configFile())
	if err != nil {
		return c, err
	}
	return c, nil
}

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
	return slugs, nil
}

func (t Track) Dirs() (map[string]struct{}, error) {
	dirs := make(map[string]struct{})

	infos, err := ioutil.ReadDir(t.path)
	if err != nil {
		return dirs, err
	}

	for _, info := range infos {
		if info.IsDir() {
			dirs[info.Name()] = struct{}{}
		}
	}
	return dirs, nil
}

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

	for problem, _ := range problems {
		_, present := dirs[problem]
		if !present {
			omissions = append(omissions, problem)
		}
	}
	return omissions, nil
}

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

	for dir, _ := range dirs {
		_, present := slugs[dir]
		if !present {
			omissions = append(omissions, dir)
		}
	}
	return omissions, nil
}

func (t Track) hasValidConfig() bool {
	_, err := t.Config()
	return err == nil
}
