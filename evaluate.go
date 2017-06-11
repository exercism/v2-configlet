package main

import (
	"fmt"

	"github.com/exercism/configlet/configlet"
)

// Evaluate sanity-checks a track.
//
// It verifies that:
// - the config is valid JSON.
// - there are no missing problems.
// - each exercise has an example file.
// - foregone problem specifications are not implemented.
// - an exercise isn't mentioned in multiple config categories.
func Evaluate(path string) bool {
	// TODO: handle this error
	track, _ := configlet.NewTrack(path)

	hasErrors := false
	if !track.HasValidConfig() {
		hasErrors = true
		fmt.Println("-> config.json is invalid. Try jsonlint.com")
	}

	configErrors := []ConfigError{
		{
			check: track.MissingProblems,
			msg:   "-> An exercise with slug '%v' is referenced in config.json, but no implementation was found.\n",
		},
		{
			check: track.UnconfiguredProblems,
			msg:   "-> An implementation for '%v' was found, but config.json does not reference this exercise.\n",
		},
		{
			check: track.ProblemsLackingExample,
			msg:   "-> The implementation for '%v' is missing an example solution.\n",
		},
		{
			check: track.ForegoneViolations,
			msg:   "-> An implementation for '%v' was found, but config.json specifies that it should be foregone (not implemented).\n",
		},
		{
			check: track.DuplicateSlugs,
			msg:   "-> The exercise '%v' was found in multiple (conflicting) categories in config.json.\n",
		},
	}

	for _, configError := range configErrors {
		result, err := configError.check()

		if err != nil {
			hasErrors = true
			fmt.Errorf("-> %v", err)
		}

		if len(result) > 0 {
			hasErrors = true
			for _, slug := range result {
				fmt.Printf(configError.msg, slug)
			}
		}
	}
	return hasErrors
}
