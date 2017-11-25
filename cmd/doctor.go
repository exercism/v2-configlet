package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var repairUUID bool

// doctorCmd represents the doctor command
var doctorCmd = &cobra.Command{
	Use:   "doctor " + pathExample,
	Short: "Check track configuration for malformed UUIDs",
	Long: `The uuid doctor command validates UUIDs within a track's configuration file.
Any malformed UUID, along with its Slug, will be captured in the diagnostic report.
To automatically fix malformed UUIDs specify the --repair flag at runtime.`,

	Run:  runDiagnostic,
	Args: cobra.ExactArgs(1),
}

func runDiagnostic(cmd *cobra.Command, args []string) {
	path := filepath.FromSlash(args[0])
	if _, err := os.Stat(path); os.IsNotExist(err) {
		ui.PrintError("path not found:", path)
		os.Exit(1)
	}

	t, err := track.New(path)
	if err != nil {
		ui.PrintError(err)
		os.Exit(1)
	}

	var badUUIDs []string
	for _, exercise := range t.Config.Exercises {
		if _, err := uuid.Parse(exercise.UUID); err == nil {
			continue
		}

		badUUIDs = append(badUUIDs, exercise.UUID)
		fmt.Fprintf(os.Stderr, "[X] %s (%s:%s)\n", exercise.UUID, t.ID, exercise.Slug)
	}

	if repairUUID {
		configPath := filepath.Join(path, "config.json")
		b, err := ioutil.ReadFile(configPath)
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}

		for _, u := range badUUIDs {
			newUUID := uuid.New().String()
			b = bytes.Replace(b, []byte(u), []byte(newUUID), 1)
		}

		err = ioutil.WriteFile(configPath, b, os.FileMode(0644))
		if err != nil {
			ui.PrintError(err)
			os.Exit(1)
		}
	}
}

func init() {
	uuidCmd.AddCommand(doctorCmd)
	doctorCmd.Flags().BoolVarP(&repairUUID, "repair", "r", false, "repair malformed UUIDs")
}
