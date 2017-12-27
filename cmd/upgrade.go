package cmd

import (
	"github.com/exercism/cli/cli"
	"github.com/exercism/cli/debug"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

// upgradeCmd downloads and installs the most recent version of Configlet.
var upgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"u"},
	Short:   "Upgrade to the latest version of Configlet.",
	Long: `Upgrade to the latest version of Configlet.

This finds and downloads the latest release, if you don't
already have it.

On Windows the old Configlet will be left on disk, marked as hidden.
The next time you upgrade, the hidden file will be overwritten.
You can always delete this file.
	`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
			debug.Verbose = verbose
		}

		cli.LatestReleaseURL = "https://api.github.com/repos/exercism/configlet/releases/latest"
		c := cli.New(Version)
		return runUpdate(c)
	},
}

// runUpdate updates Configlet to the latest available version, if it is out of date.
func runUpdate(c cli.Updater) error {
	ok, err := c.IsUpToDate()
	if err != nil {
		return err
	}

	if ok {
		ui.Print("Your Configlet version is up to date.")
		return nil
	}
	return c.Upgrade()
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolP("verbose", "v", false, "verbose output")
}
