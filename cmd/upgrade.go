package cmd

import (
	"fmt"

	"github.com/exercism/cli/cli"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

// upgradeCmd downloads and installs the most recent version of Configlet.
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to the latest version of Configlet.",
	Long: `Upgrade to the latest version of Configlet.

This finds and downloads the latest release, if you don't
already have it.

On Windows the old Configlet will be left on disk, marked as hidden.
The next time you upgrade, the hidden file will be overwritten.
You can always delete this file.
	`,

	Run: func(cmd *cobra.Command, args []string) {
		runUpdate(configletCLI)
	},
}

// runUpdate updates Configlet to the latest available version, if it is out of date.
func runUpdate(c cli.Updater) {
	ok, err := c.IsUpToDate()
	if err != nil {
		ui.PrintError(err)
	}

	if ok {
		fmt.Println("Your CLI version is up to date.")
		return
	}

	if err := c.Upgrade(); err != nil {
		ui.PrintError(err)
	}
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
}
