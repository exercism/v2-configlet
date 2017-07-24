package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
	"github.com/spf13/cobra"
)

// fmtCmd defines the fmt command
var fmtCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Format the track configuration files.",
	Long: `The fmt command formats the track's configuration files.

It ensures the following files have consistent JSON syntax and indentation:
	config.json, maintainers.json
	
It also normalizes and alphabetizes the exercise topics in the config.json file.
`,
	Run: format,
}

// formatter applies additional formatting to unmarshalled JSON files.
type formatter func(m map[string]interface{})

// formatUsageText defines how to use the fmt command
var formatUsageText = "Usage:\n  configlet fmt <path/to/track>\n"

func format(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		formatUsageFunc(cmd)
		return
	}

	path := args[0]
	var fs = []struct {
		path string
		formatter
	}{
		{
			filepath.Join(path, "config.json"),
			formatTopics,
		},
		{
			filepath.Join(path, "config", "maintainers.json"),
			nil,
		},
	}

	for _, f := range fs {
		if _, err := os.Stat(f.path); os.IsNotExist(err) {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		diff, formatted, err := formatFile(f.path, f.formatter)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if diff == "" {
			continue
		}
		err = ioutil.WriteFile(f.path, formatted, os.FileMode(0644))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		fmt.Printf("Changes made to %s:\n\n%s", f.path, diff)
	}
	return
}

func formatFile(path string, format formatter) (diff string, formatted []byte, err error) {
	f, err := os.Open(path)
	if err != nil {
		return diff, formatted, err
	}

	var m map[string]interface{}

	json.NewDecoder(f).Decode(&m)

	if format != nil {
		format(m)
	}

	original, err := ioutil.ReadFile(path)
	if err != nil {
		return diff, formatted, err
	}

	formatted, err = json.MarshalIndent(&m, "", "  ")
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
	s = strings.Replace(s, " ", "-", -1)
	return s
}

func formatUsageFunc(cmd *cobra.Command) error {
	fmt.Fprintf(os.Stderr, formatUsageText)
	return nil
}

func init() {
	RootCmd.AddCommand(fmtCmd)
	fmtCmd.SetUsageFunc(formatUsageFunc)
}
