package cmd

import (
	"fmt"
	"os"

	"github.com/exercism/configlet/track"
	"github.com/spf13/cobra"
)

var (
	// lintCmd defines the lint command.
	lintCmd = &cobra.Command{
		Use:   "lint",
		Short: "Ensure that the track is configured.",
		Long: `Verify that the config.json file is valid, and that the exercises are complete.

Call lint command with path to track:

  configlet lint <path/to/track>
`,
		Run: lint,
	}
	lintUsageText = "USAGE: configlet lint <path/to/track>\n"
)

func lint(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		lintUsageFunc(cmd)
		os.Exit(1)
	}
	var hasErrors bool
	for _, arg := range args {
		if failed := lintTrack(arg); failed {
			hasErrors = true
		}
	}
	if hasErrors {
		os.Exit(1)
	}
}

func lintTrack(path string) bool {
	t, err := track.New(path)
	if err != nil {
		fmt.Printf("-> %s\n", err)
	}

	configErrors := []struct {
		check func(track.Track) []string
		msg   string
	}{
		{
			check: missingImplementations,
			msg:   "-> An exercise with slug '%v' is referenced in config.json, but no implementation was found.\n",
		},
		{
			check: missingMetadata,
			msg:   "-> An implementation for '%v' was found, but config.json does not reference this exercise.\n",
		},
		{
			check: missingSolution,
			msg:   "-> The implementation for '%v' is missing an example solution.\n",
		},
		{
			check: missingTestSuite,
			msg:   "-> The implementation for '%v' is missing a test suite.\n",
		},
		{
			check: foregoneViolations,
			msg:   "-> An implementation for '%v' was found, but config.json specifies that it should be foregone (not implemented).\n",
		},
		{
			check: duplicateSlugs,
			msg:   "-> The exercise '%v' was found in multiple (conflicting) categories in config.json.\n",
		},
	}

	hasErrors := false
	for _, configError := range configErrors {
		slugs := configError.check(t)

		if len(slugs) > 0 {
			hasErrors = true
			for _, slug := range slugs {
				fmt.Printf(configError.msg, slug)
			}
		}
	}
	return hasErrors
}

func missingImplementations(t track.Track) []string {
	metadata := map[string]bool{}
	for _, exercise := range t.Config.Exercises {
		metadata[exercise.Slug] = false
	}
	// Don't report missing implementations on foregone exercises.
	for _, slug := range t.Config.ForegoneSlugs {
		metadata[slug] = true
	}
	for _, exercise := range t.Exercises {
		metadata[exercise.Slug] = true
	}

	slugs := []string{}
	for slug, ok := range metadata {
		if !ok {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

func missingMetadata(t track.Track) []string {
	implementations := map[string]bool{}
	for _, exercise := range t.Exercises {
		implementations[exercise.Slug] = false
	}

	// Don't report missing metadata if the exercise is deprecated or foregone.
	ignoredSlugs := append(t.Config.DeprecatedSlugs, t.Config.ForegoneSlugs...)
	for _, slug := range ignoredSlugs {
		implementations[slug] = true
	}

	for _, exercise := range t.Config.Exercises {
		implementations[exercise.Slug] = true
	}

	slugs := []string{}
	for slug, ok := range implementations {
		if !ok {
			slugs = append(slugs, slug)
		}
	}

	return slugs
}

func missingSolution(t track.Track) []string {
	solutions := map[string]bool{}
	for _, exercise := range t.Exercises {
		solutions[exercise.Slug] = exercise.IsValid()
	}
	// Don't complain about missing solutions in foregone exercises.
	for _, slug := range t.Config.ForegoneSlugs {
		solutions[slug] = true
	}

	slugs := []string{}
	for slug, ok := range solutions {
		if !ok {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

func missingTestSuite(t track.Track) []string {
	tests := map[string]bool{}
	for _, exercise := range t.Exercises {
		tests[exercise.Slug] = exercise.HasTestSuite()
	}
	// Don't complain about missing test suite in foregone exercises.
	for _, slug := range t.Config.ForegoneSlugs {
		tests[slug] = true
	}

	slugs := []string{}
	for slug, ok := range tests {
		if !ok {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

func foregoneViolations(t track.Track) []string {
	violations := map[string]bool{}
	for _, slug := range t.Config.ForegoneSlugs {
		violations[slug] = true
	}

	slugs := []string{}
	for _, exercise := range t.Exercises {
		if violations[exercise.Slug] {
			slugs = append(slugs, exercise.Slug)
		}
	}

	return slugs
}

func duplicateSlugs(t track.Track) []string {
	counts := map[string]int{}
	for _, slug := range t.Config.ForegoneSlugs {
		counts[slug]++
	}
	for _, slug := range t.Config.DeprecatedSlugs {
		counts[slug]++
	}
	for _, exercise := range t.Config.Exercises {
		counts[exercise.Slug]++
	}

	slugs := []string{}
	for slug, count := range counts {
		if count > 1 {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

func lintUsageFunc(cmd *cobra.Command) error {
	fmt.Fprintf(os.Stderr, lintUsageText)
	return nil
}

func init() {
	RootCmd.AddCommand(lintCmd)
	lintCmd.SetUsageFunc(lintUsageFunc)
}
