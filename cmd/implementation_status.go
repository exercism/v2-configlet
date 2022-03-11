package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

var pathExampleProblemSpecifications = "<path/to/problemspecifications>"

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

const (
	implementationStatusOK = iota
	implementationStatusDeprecated
	implementationStatusForegone
)

// listProblemSpecs gives a map of exercise name and whether the exercise is deprecated.
func listProblemSpecs(path string) (map[string]int, error) {
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]int)
	for _, f := range fis {
		_, err := os.Stat(filepath.Join(path, f.Name(), ".deprecated"))
		if os.IsNotExist(err) {
			ret[f.Name()] = implementationStatusOK
		} else {
			ret[f.Name()] = implementationStatusDeprecated
		}
	}

	return ret, nil
}

func trackConfigExercises(track string) (map[string]int, error) {
	var trackConfig struct {
		Foregone  []string `json:"foregone"`
		Exercises []struct {
			Slug       string `json:"slug"`
			Deprecated bool   `json:"deprecated"`
		} `json:"exercises"`
	}

	f, err := os.Open(filepath.Join(track, "config.json"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&trackConfig); err != nil {
		return nil, err
	}

	names := make(map[string]int)
	for _, x := range trackConfig.Exercises {
		if x.Deprecated {
			names[x.Slug] = implementationStatusDeprecated
		} else {
			names[x.Slug] = implementationStatusOK
		}
	}
	for _, f := range trackConfig.Foregone {
		names[f] = implementationStatusForegone
	}
	return names, nil
}

func runImplementationStatus(problemSpecifications, track string) {

	trackConfig, err := trackConfigExercises(track)
	if err != nil {
		ui.PrintError(err.Error())
		return
	}

	namesPS, err := listProblemSpecs(filepath.Join(problemSpecifications, "exercises"))
	if err != nil {
		ui.PrintError(err.Error())
		return
	}

	results := make([]string, 0)

	for k, v := range namesPS {
		if _, ok := trackConfig[k]; !ok {
			results = append(results, fmt.Sprintf("The exercise with slug '%s' is not implemented in this track.", k))
		} else {
			if v == implementationStatusDeprecated && trackConfig[k] == implementationStatusOK {
				results = append(results, fmt.Sprintf("The exercise with slug '%s' exists in this track, but has been deprecated in the specifications repository.", k))
			}
		}
	}

	for k := range trackConfig {
		if _, ok := namesPS[k]; !ok {
			results = append(results, fmt.Sprintf("The exercise with slug '%s' exists in this track, but has been removed from the specifications repository.", k))
		}
	}

	sort.Strings(results)
	for _, r := range results {
		ui.Print(r)
	}
}

func init() {
	RootCmd.AddCommand(implementationStatusCmd)
}
