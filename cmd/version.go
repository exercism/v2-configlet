package cmd

import (
	"fmt"

	"github.com/exercism/configlet/ui"
	"github.com/nywilken/cli/cli"
	"github.com/spf13/cobra"
)

// Version is the current version of the tool.
const Version = "3.6.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Output the current version of the tool",
	Long:    "Output the current version of the tool",
	Example: fmt.Sprintf("  %s version", binaryName),
	Run:     runVersion,
	Args:    cobra.ExactArgs(0),
}

func runVersion(cmd *cobra.Command, args []string) {
	cli.ReleaseURL = "https://api.github.com/repos/exercism/configlet/releases"
	// we don't want any UI formatting prepended to this
	fmt.Printf("%s version %s\n", binaryName, Version)

	c := cli.New(Version)
	ok, err := c.IsUpToDate()
	if err != nil {
		ui.PrintError(err)
	}

	if !ok {
		msg := fmt.Sprintf("There is a newer version of Configlet available (%s)", c.LatestRelease.Version())
		ui.Print(msg)
	}
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
