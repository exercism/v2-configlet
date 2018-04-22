package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/exercism/configlet/ui"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

// verbose flag for fmt command.
var fmtVerbose bool

// fmtCmd defines the fmt command
var fmtCmd = &cobra.Command{
	Use:   "fmt " + pathExample,
	Short: "Format the track configuration files",
	Long: `The fmt command formats the track's configuration files.

It ensures the following files have consistent JSON syntax and indentation:
	config.json, maintainers.json

It also normalizes and alphabetizes the exercise topics in the config.json file.
`,
	Example: fmt.Sprintf("  %s fmt %s --verbose", binaryName, pathExample),
	Run:     runFmt,
	Args:    cobra.ExactArgs(1),
}

// formatter applies additional formatting to unmarshalled JSON files.
type formatter func(m map[string]interface{})

// orderer applies an ordering to unmarshalled JSON files
type orderer func(map[string]interface{}) OrderedMap

func runFmt(cmd *cobra.Command, args []string) {
	path := args[0]
	var fs = []struct {
		path string
		formatter
		orderer
	}{
		{
			filepath.Join(path, "config.json"),
			formatTopics,
			orderConfig,
		},
		{
			filepath.Join(path, "config", "maintainers.json"),
			nil,
			nil,
		},
	}

	var changes string

	for _, f := range fs {
		if _, err := os.Stat(f.path); os.IsNotExist(err) {
			ui.PrintError("path not found:", f.path)
			os.Exit(1)
		}
		diff, formatted, err := formatFile(f.path, f.formatter, f.orderer)
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}
		if diff == "" {
			continue
		}
		err = ioutil.WriteFile(f.path, formatted, os.FileMode(0644))
		if err != nil {
			ui.PrintError(err.Error())
			continue
		}
		if fmtVerbose {
			ui.Print(f.path, "\n\n", diff)
		}
		changes += fmt.Sprintf("%s\n", f.path)
	}

	if changes != "" {
		ui.Print("changes made to:\n", changes)
	}
}

func formatFile(path string, format formatter, order orderer) (diff string, formatted []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return diff, formatted, err
	}

	var m map[string]interface{}
	if err = json.NewDecoder(f).Decode(&m); err != nil {
		return diff, formatted, err
	}

	if format != nil {
		format(m)
	}

	var om interface{}
	if order == nil {
		om = m
	} else {
		om = order(m)
	}

	original, err := ioutil.ReadFile(path)
	if err != nil {
		return diff, formatted, err
	}

	formatted, err = json.MarshalIndent(&om, "", "  ")
	if err != nil {
		return diff, formatted, err
	}

	diff, err = difflib.GetUnifiedDiffString(
		difflib.UnifiedDiff{
			A: difflib.SplitLines(string(original)),
			B: difflib.SplitLines(string(formatted)),
		})

	return diff, formatted, err
}

func formatTopics(m map[string]interface{}) {
	exercises, ok := m["exercises"].([]interface{})
	if !ok {
		return
	}
	for _, e := range exercises {
		exercise, ok := e.(map[string]interface{})
		if !ok {
			continue
		}
		topics, ok := exercise["topics"].([]interface{})
		if !ok {
			continue
		}
		var sorted []string
		for _, t := range topics {
			topic, ok := t.(string)
			if !ok {
				continue
			}
			sorted = append(sorted, normaliseTopic(topic))
		}
		sort.Strings(sorted)
		exercise["topics"] = sorted
	}
}

func normaliseTopic(t string) string {
	s := strings.ToLower(t)

	// we only want to let through letters and underscores
	// for the final output.
	// hyphens and whitespace are allowed through for now
	// as word delimiters to be replaced with underscores.
	reg := regexp.MustCompile(`[^a-z\s-_]+`)
	s = reg.ReplaceAllString(s, "")

	// output to be snake_case
	reg = regexp.MustCompile(`[\s-]+`)
	s = reg.ReplaceAllString(s, "_")
	return s
}

// orderConfig applies an ordering to config.json
//
// This ordering was specified in
// https://github.com/exercism/meta/issues/95#issuecomment-341991512
func orderConfig(m map[string]interface{}) OrderedMap {
	exercises, ok := m["exercises"].([]interface{})
	if ok {
		var orderedExercises []OrderedMap
		for _, e := range exercises {
			exercise, ok := e.(map[string]interface{})
			if !ok {
				continue
			}
			orderedExercises = append(orderedExercises, WithOrdering(exercise,
				"slug",
				"uuid",
				"core",
				"unlocked_by",
				"difficulty",
				"topics",
			))
		}
		m["exercises"] = orderedExercises
	}
	return WithOrdering(m,
		"track_id",
		"language",
		"active",
		"blurb",
		"gitter",
		"checklist_issue",
		"ignore_pattern",
		"solution_pattern",
		"test_pattern",
		"foregone",
		"exercises",
	)
}

func init() {
	RootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().BoolVarP(&fmtVerbose, "verbose", "v", false, "display the diff of the formatted changes.")
}
