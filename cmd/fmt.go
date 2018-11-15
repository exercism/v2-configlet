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

var (
	// verbose flag for fmt command.
	fmtVerbose bool

	// test flag for fmt command displays the proposed changes
	fmtTest bool
)

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

		if diffFound, err := runFmt(args[0], args[0]); err != nil {
			ui.PrintError(err.Error())
			os.Exit(1)
		} else if diffFound && fmtTest {
			os.Exit(2)
		}

	},

	Args: cobra.ExactArgs(1),
}

func runFmt(inDir, outDir string) (bool, error) {
	if _, err := os.Stat(filepath.Join(outDir, "config")); os.IsNotExist(err) {
		os.Mkdir(filepath.Join(outDir, "config"), os.ModePerm)
	}

	if fmtTest {
		fmtVerbose = true
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
			if fmtVerbose {
				ui.Print(fmt.Sprintf("%s\n\n%s", f.inPath, diff))
			}
			changes += fmt.Sprintf("%s\n", f.inPath)

		}
	}
	diffFound := changes != ""
	if diffFound {
		if fmtTest {
			ui.Print("no changes were made to:\n", changes)
		} else {
			ui.Print("changes made to:\n", changes)
		}
	}
	if err := errs.ErrorOrNil(); err != nil {
		return diffFound, err
	}
	return diffFound, nil
}

func formatFile(cfg ConfigSerializer, inPath, outPath string) (string, error) {
	src, err := ioutil.ReadFile(inPath)
	if err != nil {
		return "", err
	}
	if err := cfg.LoadFromFile(inPath); err != nil {
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
	if diff != "" && !fmtTest {
		if err := ioutil.WriteFile(outPath, dst, os.FileMode(0644)); err != nil {
			return "", err
		}
	}
	return diff, nil
}

func init() {
	RootCmd.AddCommand(fmtCmd)
	fmtCmd.Flags().BoolVarP(&fmtVerbose, "verbose", "v", false, "display the diff of the formatted changes.")
	fmtCmd.Flags().BoolVarP(&fmtTest, "test", "t", false, "display the proposed changes, but do not make them.")
}
