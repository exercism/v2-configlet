package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is the current version of the tool.
const Version = "3.4.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Output the current version of the tool",
	Long:    "Output the current version of the tool",
	Example: fmt.Sprintf("  %s version", binaryName),
	Run:     runVersion,
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Printf("%s version %s\n", binaryName, Version)
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
