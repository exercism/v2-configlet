package main

import (
	"fmt"

	"github.com/exercism/configlet/configlet"
)

// Evaluate sanity-checks a track.
// It verifies that the config is valid JSON, that
// there are no missing problems, that each problem
// has an example file, and that foregone problems
// are not present. It also checks that a problem
// isn't mentioned in multiple categories.
func Evaluate(path string) bool {
	track := configlet.NewTrack(path)

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
