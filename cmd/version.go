package cmd

import (
	"fmt"

	"github.com/exercism/cli/cli"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

// Version is the current version of the tool.
const Version = "3.7.0"

var checkLatestVersion bool

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
	fmt.Printf("configlet version %s\n", Version)

	if !checkLatestVersion {
		return
	}

	c := cli.New(Version)
	ok, err := c.IsUpToDate()
	if err != nil {
		ui.PrintError(err)
		return
	}
	msg := "Your version is up to date."

	if !ok {
		msg = fmt.Sprintf("A new version is available. Run `%s upgrade` to update to %s", binaryName, c.LatestRelease.Version())
	}
	fmt.Println(msg)
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&checkLatestVersion, "latest", "l", false, "check latest available version.")
}
