package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/exercism/cli/cli"
	"github.com/exercism/configlet/ui"
	"github.com/spf13/cobra"
)

var (
	// GoReleaser injects this variable using ldflags
	version = "dev"
	// binaryName is this tool's given name. While we have named it configlet,
	// others may choose to rename it. This var enables the use of it's name,
	// whatever it is, in any help/usage text.
	binaryName = os.Args[0]
	//configletCLI is an updatable CLI binary.
	configletCLI *cli.CLI
	// pathExample illustrates the path argument necessary for some commands.
	pathExample = "<path/to/track>"
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     binaryName,
	Short:   "A tool for managing Exercism language track repositories.",
	Long:    binaryName + " version " + version + "\n\n" + "A tool for managing Exercism language track repositories.",
	Example: rootExampleText(),
}

func rootExampleText() string {
	cmds := []string{
		"%[1]s fmt %[2]s",
		"%[1]s generate %[2]s",
		"%[1]s lint %[2]s",
	}
	s := "  " + strings.Join(cmds, "\n\n  ")
	return fmt.Sprintf(s, binaryName, pathExample)
}

// Execute adds all child commands to the root command & sets flags
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		ui.PrintError(err.Error())
		os.Exit(-1)
	}
}

func init() {
	cli.ReleaseURL = "https://api.github.com/repos/exercism/configlet/releases"
	configletCLI = cli.New(version)
}
