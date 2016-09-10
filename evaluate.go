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
// - each problem has an example file.
// - foregone problems are not implemented.
// - a problem isn't mentioned in multiple config categories.
func Evaluate(path string) bool {
	// TODO: handle this error
	track, _ := configlet.NewTrack(path)

	hasErrors := false
	if !track.HasValidConfig() {
		hasErrors = true
		fmt.Println("-> config.json is invalid. Try jsonlint.com")
	}

	configErrors := []ConfigError{
		ConfigError{
			check: track.MissingProblems,
			msg:   "-> No directory found for %v.\n",
		},
		ConfigError{
			check: track.UnconfiguredProblems,
			msg:   "-> config.json does not include %v.\n",
		},
		ConfigError{
			check: track.ProblemsLackingExample,
			msg:   "-> missing example solution in %v.\n",
		},
		ConfigError{
			check: track.ForegoneViolations,
			msg:   "-> %v should not be implemented.\n",
		},
		ConfigError{
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
