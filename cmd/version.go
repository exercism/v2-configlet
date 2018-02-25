package cmd

import (
	"fmt"

	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

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
	// we don't want any UI formatting prepended to this
	fmt.Printf("configlet version %s\n", configletCLI.Version)

	if !checkLatestVersion {
		return
	}

	ok, err := configletCLI.IsUpToDate()
	if err != nil {
		ui.PrintError(err)
		return
	}
	msg := "Your CLI version is up to date."

	if !ok {
		msg = fmt.Sprintf("A new CLI version is available. Run `%s upgrade` to update to %s", binaryName, configletCLI.LatestRelease.Version())
	}
	fmt.Println(msg)
}

func init() {
	RootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVarP(&checkLatestVersion, "latest", "l", false, "check latest available version.")
}
