package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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

func (c Config) IgnoredDirs() []string {
	return append(c.Ignored, ".git", "bin")
}
