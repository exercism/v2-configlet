package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
	multierror "github.com/hashicorp/go-multierror"
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

func runFmt(inDir, outDir string, verbose bool) error {
	if _, err := os.Stat(filepath.Join(outDir, "config")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(outDir, "config"), os.ModePerm)
	}

	var fs = []struct {
		inPath  string
		outPath string
		cfg     ConfigSerializer
	}{
		{
			filepath.Join(inDir, "config.json"),
			filepath.Join(outDir, "config.json"),
			&track.Config{},
		},
		{
			filepath.Join(inDir, "config", "maintainers.json"),
			filepath.Join(outDir, "config", "maintainers.json"),
			&track.MaintainerConfig{},
		},
	}

	var changes string

	errs := &multierror.Error{}
	for _, f := range fs {
		diff, err := formatFile(f.cfg, f.inPath, f.outPath)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}
		if diff != "" {
			if verbose {
				ui.Print(fmt.Sprintf("%s\n\n%s", f.inPath, diff))
			}
			changes += fmt.Sprintf("%s\n", f.inPath)
		}
	}
	if changes != "" {
		ui.Print("changes made to:\n", changes)
	}
	if err := errs.ErrorOrNil(); err != nil {
		return err
	}
	return nil
}

func formatFile(cfg ConfigSerializer, inPath, outPath string) (string, error) {
	src, err := ioutil.ReadFile(inPath)
	if err != nil {
		return "", err
	}
	if err := cfg.NewConfigFromFile(inPath); err != nil {
		return "", err
	}
	dst, err := cfg.ToJSON()
	if err != nil {
		return "", err
	}
	dst = []byte(string(fmt.Sprintf("%s\n", dst)))

	a := difflib.SplitLines(string(src))
	b := difflib.SplitLines(string(dst))
	diff, err := difflib.GetUnifiedDiffString(difflib.UnifiedDiff{A: a, B: b})
	if err != nil {
		return "", err
	}
	if diff != "" {
		if err := ioutil.WriteFile(outPath, dst, os.FileMode(0644)); err != nil {
			return "", err
		}
	}
	return diff, nil
}

func init() {
	RootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().BoolVarP(&fmtVerbose, "verbose", "v", false, "display the diff of the formatted changes.")
}
