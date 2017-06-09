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
			msg:   "-> No directory found for %v.\n",
		},
		{
			check: track.UnconfiguredProblems,
			msg:   "-> config.json does not include %v.\n",
		},
		{
			check: track.ProblemsLackingExample,
			msg:   "-> missing example solution in %v.\n",
		},
		{
			check: track.ForegoneViolations,
			msg:   "-> %v should not be implemented.\n",
		},
		{
			check: track.DuplicateSlugs,
			msg:   "-> %v found in multiple categories.\n",
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
			fmt.Printf(configError.msg, result)
		}
	}
	return hasErrors
}
