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
	Run: func(cmd *cobra.Command, args []string) {
		if err := runFmt(args[0], args[0], fmtVerbose); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		}
	},

	Args: cobra.ExactArgs(1),
}

// formatter applies additional formatting to unmarshalled JSON files.
type formatter func(m map[string]interface{})

// orderer applies an ordering to unmarshalled JSON files
type orderer func(map[string]interface{}) OrderedMap

func runFmt(inDir, outDir string, verbose bool) error {
	var fs = []struct {
		inPath  string
		outPath string
		formatter
		orderer
	}{
		{
			filepath.Join(inDir, "config.json"),
			filepath.Join(outDir, "config.json"),
			formatTopics,
			orderConfig,
		},
		{
			filepath.Join(inDir, "config", "maintainers.json"),
			filepath.Join(outDir, "config", "maintainers.json"),
			formatMaintainers,
			orderMaintainers,
		},
	}

	// This is for the tests.
	// It will go away in a subsequent refactoring.
	os.Mkdir(filepath.Join(outDir, "config"), os.ModePerm)

	var changes string

	for _, f := range fs {
		diff, err := formatFile(f.inPath, f.outPath, f.formatter, f.orderer)
		if err != nil {
			return err
		}
		if diff == "" {
			continue
		}
		if verbose {
			ui.Print(fmt.Sprintf("%s\n\n%s", f.inPath, diff))
		}
		changes += fmt.Sprintf("%s\n", f.inPath)
	}

	if changes != "" {
		ui.Print("changes made to:\n", changes)
	}
	return nil
}

func formatFile(inPath, outPath string, format formatter, order orderer) (string, error) {
	if _, err := os.Stat(inPath); os.IsNotExist(err) {
		return "", fmt.Errorf("path not found: %s", inPath)
	}

	f, err := os.Open(inPath)
	if err != nil {
		return "", err
	}

	var m map[string]interface{}
	if err = json.NewDecoder(f).Decode(&m); err != nil {
		return "", err
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

	original, err := ioutil.ReadFile(inPath)
	if err != nil {
		return "", err
	}

	formatted, err := json.MarshalIndent(&om, "", "  ")
	if err != nil {
		return "", err
	}

	src := difflib.SplitLines(strings.TrimSuffix(string(original), "\n"))
	dst := difflib.SplitLines(string(formatted))
	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{A: src, B: dst})
	if diff == "" || err != nil {
		return "", err
	}

	formatted = []byte(fmt.Sprintf("%s\n", formatted))
	err = ioutil.WriteFile(outPath, formatted, os.FileMode(0644))
	return diff, err
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
			// Topics are null.
			continue
		}
		// Ensure that topics are an empty list rather than null.
		sorted := make([]string, 0, len(topics))
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

// orderMaintainers applies an ordering to maintainers.json
//
// This ordering was specified in
// https://github.com/exercism/meta/issues/95#issuecomment-341991512
func orderMaintainers(m map[string]interface{}) OrderedMap {
	maintainers, ok := m["maintainers"].([]interface{})
	orderedMaintainers := make([]OrderedMap, 0, len(maintainers))
	if ok {
		for _, m := range maintainers {
			maintainer, ok := m.(map[string]interface{})
			if !ok {
				continue
			}
			orderedMaintainers = append(orderedMaintainers, WithOrdering(maintainer,
				"github_username",
				"alumnus",
				"show_on_website",
				"name",
				"link_text",
				"link_url",
				"avatar_url",
				"bio",
			))
		}
	}
	m["maintainers"] = orderedMaintainers
	return WithOrdering(m,
		"docs_url",
		"maintainers",
	)
}

// Add missing fields
func formatMaintainers(m map[string]interface{}) {
	maintainers, ok := m["maintainers"].([]interface{})
	if !ok {
		return
	}
	for _, m := range maintainers {
		maintainer, ok := m.(map[string]interface{})
		if !ok {
			continue
		}

		maintainerKeys := []string{
			"github_username",
			"alumnus",
			"show_on_website",
			"name",
			"link_text",
			"link_url",
			"avatar_url",
			"bio",
		}
		for _, key := range maintainerKeys {
			_, ok := maintainer[key].(interface{})
			if !ok {
				maintainer[key] = nil
			}
		}
	}
}

func init() {
	RootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().BoolVarP(&fmtVerbose, "verbose", "v", false, "display the diff of the formatted changes.")
}
