package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

var pathExampleProblemSpecifications = "<path/to/problem-specifications>"

// implementationStatusCmd gives information about exercises that should be
// worked on (e.g., new version or not implemented yet).
var implementationStatusCmd = &cobra.Command{
	Use: "implementation-status",

	Short: "Show the status of implementation for any exercise that needs work.",
	Long: `Show the status of implementation for any exercise that needs work.

For example, if upstream prepares a new exercise it is important to either
implement it in a track or declare it foregone. If a new version is published the update
should be made available soon.
`,

	Args: cobra.ExactArgs(2),

	Example: fmt.Sprintf("  %s implementation-status %s %s", binaryName, pathExampleProblemSpecifications, pathExample),
	Run: func(cmd *cobra.Command, args []string) {
		runImplementationStatus(args[0], args[1])
	},
}

// listProblemSpecs gives a map of exercise name and whether the exercise is deprecated.
func listProblemSpecs(path string) (map[string]bool, error) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]bool)
	for _, f := range fis {
		_, err := os.Stat(filepath.Join(path, f.Name(), ".deprecated"))
		if os.IsNotExist(err) {
			ret[f.Name()] = true
		} else if err == nil {
			ret[f.Name()] = false
		} else {
			return nil, err
		}
	}

	return ret, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func runImplementationStatus(specPath, trackPath string) {

	trackConfigPath := filepath.Join(trackPath, "config.json")
	track.ProblemSpecificationsPath = filepath.Join(filepath.Dir(trackConfigPath), track.ProblemSpecificationsDir)
	if specPath != "" {
		track.ProblemSpecificationsPath = specPath
	}

	if _, err := os.Stat(track.ProblemSpecificationsPath); os.IsNotExist(err) {
		ui.PrintError("path not found:", track.ProblemSpecificationsPath)
		os.Exit(1)
	}

	trackConfig, err := track.NewConfig(trackConfigPath)
	if err != nil {
		ui.PrintError(err.Error())
		return
	}

	namesPS, err := listProblemSpecs(filepath.Join(track.ProblemSpecificationsPath, "exercises"))
	if err != nil {
		ui.PrintError(err.Error())
		return
	}

	trackExerciseSlugs := make(map[string]bool)

	for _, exercise := range trackConfig.Exercises {
		trackExerciseSlugs[exercise.Slug] = true
	}

	for slug, active := range namesPS {
		if trackExerciseSlugs[slug] {
			if stringInSlice(slug, trackConfig.DeprecatedSlugs) || !active {
				ui.Print(fmt.Sprintf("The exercise '%s' exists in this track but is deprecated.", slug))
			} else if stringInSlice(slug, trackConfig.ForegoneSlugs) {
				ui.Print(fmt.Sprintf("The exercise '%s' exists in this track but is forgone.", slug))
			}
		} else {
			if !stringInSlice(slug, trackConfig.ForegoneSlugs) {
				ui.Print(fmt.Sprintf("The exercise '%s' does not exist in this track.", slug))
			}
		}
	}

	for slug := range trackExerciseSlugs {
		if _, ok := namesPS[slug]; !ok {
			ui.Print(fmt.Sprintf("The exercise '%s' exists in the track but is not in the problem-specifications repository.", slug))
		}
	}
}

func init() {
	RootCmd.AddCommand(implementationStatusCmd)
}
