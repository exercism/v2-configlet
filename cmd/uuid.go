package cmd

import (
	"fmt"

	"github.com/mattetti/uuid"
	"github.com/spf13/cobra"
)

// uuidCmd represents the uuid command
var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate a UUID",
	Long: `Generate a UUID.

Each Exercism exercise needs a unique and unmutable UUID. This needs to be
unique across the entire platform, not just within a track. In other words, if
you have 'clock' in Go and 'clock' in Haskell, they need to have different
UUIDs, even though they are based on the same problem specification.
`,
	Example: fmt.Sprintf("  %s uuid", binaryName),
	Run:     runUUID,
	Args:    cobra.ExactArgs(0),
}

// runUUID prints out a unique exercise UUID
func runUUID(cmd *cobra.Command, args []string) {
	fmt.Println(uuid.GenUUID())
	return
}

func init() {
	RootCmd.AddCommand(uuidCmd)
}
