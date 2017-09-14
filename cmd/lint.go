package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

// UUIDValidationURL is the endpoint to Exercism's UUID validation service.
var UUIDValidationURL = "http://exercism.io/api/v1/uuids"

// disableHTTPChecks flag skips HTTP based checks when passed.
var disableHTTPChecks bool

// lintCmd defines the lint command.
var lintCmd = &cobra.Command{
	Use:   "lint " + pathExample,
	Short: "Ensure that the track is configured correctly",
	Long: `The lint command checks for any discrepancies in a track's configuration files.

It ensures the following files are valid JSON:
	config.json, maintainers.json

It also checks that the exercises defined in the config.json file are complete.
`,
	Example: fmt.Sprintf("  %s lint %s", binaryName, pathExample),
	Run:     runLint,
	Args:    cobra.ExactArgs(1),
}

func runLint(cmd *cobra.Command, args []string) {
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
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ui.PrintError("path not found:", path)
		os.Exit(1)
	}
	t, err := track.New(path)
	if err != nil {
		ui.PrintError(err.Error())
	}

	configErrors := []struct {
		check func(track.Track) []string
		msg   string
	}{
		{
			check: missingImplementations,
			msg:   "An exercise with slug '%v' is referenced in config.json, but no implementation was found.",
		},
		{
			check: missingMetadata,
			msg:   "An implementation for '%v' was found, but config.json does not reference this exercise.",
		},
		{
			check: missingSolution,
			msg:   "The implementation for '%v' is missing an example solution.",
		},
		{
			check: missingTestSuite,
			msg:   "The implementation for '%v' is missing a test suite.",
		},
		{
			check: missingUUID,
			msg:   "The exercise '%v' was found in config.json, but does not have a UUID.",
		},
		{
			check: foregoneViolations,
			msg:   "An implementation for '%v' was found, but config.json specifies that it should be foregone (not implemented).",
		},
		{
			check: duplicateSlugs,
			msg:   "The exercise '%v' was found in multiple (conflicting) categories in config.json.",
		},
		{
			check: duplicateUUID,
			msg:   "The following UUID occurs multiple times. Each exercise UUID must be unique.\n%v",
		},
		{
			check: duplicateTrackUUID,
			msg:   "The following UUID was found in multiple Exercism tracks. Each exercise UUID must be unique across tracks.\n%v",
		},
		{
			check: invalidRegexPatterns,
			msg:   "-> The following pattern %q failed to compile. Please check the Regex pattern.\n",
		},
	}

	hasErrors := false
	for _, configError := range configErrors {
		failedItems := configError.check(t)

		if len(failedItems) > 0 {
			hasErrors = true
			for _, item := range failedItems {
				ui.Print(fmt.Sprintf(configError.msg, item))

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

func missingUUID(t track.Track) []string {
	slugs := []string{}
	for _, exercise := range t.Config.Exercises {
		if exercise.UUID == "" {
			slugs = append(slugs, exercise.Slug)
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

func duplicateUUID(t track.Track) []string {
	uuids := []string{}
	seen := map[string]bool{}
	for _, exercise := range t.Config.Exercises {
		if exercise.UUID == "" {
			continue
		}

		if seen[exercise.UUID] {
			uuids = append(uuids, exercise.UUID)
		}

		seen[exercise.UUID] = true
	}

	return uuids
}

func duplicateTrackUUID(t track.Track) []string {
	if disableHTTPChecks {
		return []string{}
	}

	// Build up set of uuids to validate.
	uuids := []string{}
	for _, exercise := range t.Config.Exercises {
		if exercise.UUID == "" {
			continue
		}
		uuids = append(uuids, exercise.UUID)
	}

	payload := struct {
		TrackID string   `json:"track_id"`
		UUIDs   []string `json:"uuids"`
	}{
		TrackID: t.ID,
		UUIDs:   uuids,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}

	resp, err := http.Post(UUIDValidationURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		ui.PrintError(err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		result := struct{ UUIDs []string }{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}

		return result.UUIDs
	}

	return []string{}
}

func invalidRegexPatterns(t track.Track) []string {
	patterns := t.Config.Patterns()

	failedPatterns := []string{}
	for _, pattern := range patterns {
		if _, err := regexp.Compile(pattern); err != nil {
			failedPatterns = append(failedPatterns, regexp.QuoteMeta(pattern))
		}
	}

	return failedPatterns
}

func invalidRubyRegexPatterns(t track.Track) []string {
	if disableHTTPChecks {
		return []string{}
	}

	body, err := json.Marshal(t.Config.PatternGroup)
	if err != nil {
		fmt.Fprintf(os.Stderr, "-> %s\n", err.Error())
		os.Exit(1)
	}

	resp, err := http.Post(UUIDValidationURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		fmt.Fprintf(os.Stderr, "-> %s\n", err.Error())
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnprocessableEntity {
		result := struct{ Patterns []string }{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			fmt.Fprintf(os.Stderr, "-> %s\n", err.Error())
			os.Exit(1)
		}

		return result.Patterns
	}

	return []string{}
}

func init() {
	RootCmd.AddCommand(lintCmd)
	lintCmd.Flags().BoolVar(&disableHTTPChecks, "no-http", false, "Disable remote HTTP-based linting.")
}
