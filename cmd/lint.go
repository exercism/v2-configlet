package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

var (
	// UUIDValidationURL is the endpoint to Exercism's UUID validation service.
	UUIDValidationURL = "http://exercism.io/api/v1/uuids"
	// noHTTP flag indicates if HTTP-based lint checks have been disabled at runtime.
	noHTTP bool
	// trackID flag allows the user to specify the ID of the track,
	// for example if it is different to the local directory name
	trackID string
)

// lintCmd defines the lint command.
var lintCmd = &cobra.Command{
	Use:   "lint " + pathExample,
	Short: "Ensure that the track is configured correctly",
	Long: `The lint command checks for any discrepancies in a track's configuration files.

It ensures the following files are valid JSON:
	config.json, maintainers.json

It also checks that the exercises defined in the config.json file are complete.
`,
	Example: lintExampleText(),
	Run:     runLint,
	Args:    cobra.ExactArgs(1),
}

func lintExampleText() string {
	cmds := []string{
		"%[1]s lint %[2]s",
		"%[1]s lint %[2]s --no-http",
		"%[1]s lint %[2]s --track-id=<track id>",
	}
	s := "  " + strings.Join(cmds, "\n\n  ")
	return fmt.Sprintf(s, binaryName, pathExample)
}

func runLint(cmd *cobra.Command, args []string) {
	var hasErrors bool
	for _, arg := range args {
		hasErrors = lintTrack(arg) || hasErrors
	}
	if hasErrors {
		os.Exit(1)
	}
}

func lintTrack(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ui.PrintError("path not found:", path)
		return true
	}

	t, err := track.New(path)
	if err != nil {
		ui.PrintError(err.Error())
		return true
	}

	if trackID != "" {
		t.ID = trackID
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
	}

	var hasErrors bool
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

type folder struct{ add, remove []string }

func (f folder) fold(m map[string]int) map[string]int {
	for _, slug := range f.add {
		m[slug]++
	}
	for _, slug := range f.remove {
		delete(m, slug)
	}
	return m
}

func asFolder(add, remove []string) folder { return folder{add: add, remove: remove} }
func reversed(f folder) folder             { f.add, f.remove = f.remove, f.add; return f }

func applyFolders(fs ...folder) map[string]int {
	m := map[string]int{}
	for _, f := range fs {
		m = f.fold(m)
	}
	return m
}

func foldSlugsByExistence(fs ...folder) []string {
	return foldSlugsByCount(0, fs...)
}

func foldSlugsByCount(threshold int, fs ...folder) []string {
	var slugs []string
	for slug, count := range applyFolders(fs...) {
		if count > threshold {
			slugs = append(slugs, slug)
		}
	}
	return slugs
}

func missingImplementations(t track.Track) []string {
	return foldSlugsByExistence(
		folder{add: t.Config.ExerciseSlugs()},
		folder{remove: t.Config.ForegoneSlugs},
		folder{remove: t.ExerciseSlugs()})
}

func missingMetadata(t track.Track) []string {
	return foldSlugsByExistence(
		folder{add: t.ExerciseSlugs()},
		folder{remove: t.Config.DeprecatedSlugs},
		folder{remove: t.Config.ForegoneSlugs},
		folder{remove: t.Config.ExerciseSlugs()})
}

func missingSolution(t track.Track) []string {
	return foldSlugsByExistence(
		reversed(asFolder(track.Exercises(t.Exercises).Fold(track.Exercise.IsValid))),
		folder{remove: t.Config.ForegoneSlugs})
}

func missingTestSuite(t track.Track) []string {
	return foldSlugsByExistence(
		reversed(asFolder(track.Exercises(t.Exercises).Fold(track.Exercise.HasTestSuite))),
		folder{remove: t.Config.ForegoneSlugs})
}

func missingUUID(t track.Track) []string {
	return foldSlugsByExistence(
		reversed(asFolder(track.ExerciseMetadataList(t.Config.Exercises).
			Fold(track.ExerciseMetadata.HasUUID))))
}

func foregoneViolations(t track.Track) []string {
	return foldSlugsByCount(1,
		folder{add: t.Config.ForegoneSlugs},
		folder{add: t.ExerciseSlugs()})
}

func duplicateSlugs(t track.Track) []string {
	return foldSlugsByCount(1,
		folder{add: t.Config.ForegoneSlugs},
		folder{add: t.Config.DeprecatedSlugs},
		folder{add: t.Config.ExerciseSlugs()})
}

func duplicateUUID(t track.Track) []string {
	return foldSlugsByCount(1,
		folder{add: t.Config.ExerciseUUIDs(false)})
}

func duplicateTrackUUID(t track.Track) []string {
	if noHTTP {
		return []string{}
	}

	payload := struct {
		TrackID string   `json:"track_id"`
		UUIDs   []string `json:"uuids"`
	}{
		TrackID: t.ID,
		UUIDs:   t.Config.ExerciseUUIDs(false),
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

	return nil
}

func init() {
	RootCmd.AddCommand(lintCmd)
	lintCmd.Flags().BoolVar(&noHTTP, "no-http", false, "Disable remote HTTP-based linting.")
	lintCmd.Flags().StringVar(&trackID, "track-id", "", "Specify the track ID (defaults to local directory name).")
}
