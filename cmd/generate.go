package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/exercism/configlet/track"
	"github.com/spf13/cobra"
)

var (
	genSlug string
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate exercise READMEs for an Exercism language track",
	Long: `Generate READMEs for Exercism exercises based on contents of
	a number of different files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "USAGE: configlet generate path/to/track")
			os.Exit(1)
		}

		path, err := filepath.Abs(filepath.FromSlash(args[0]))
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		root := filepath.Dir(path)
		trackID := filepath.Base(path)

		var exercises []track.Exercise

		if genSlug != "" {
			exercises = append(exercises, track.Exercise{Slug: genSlug})
		} else {
			track, err := track.New(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, err.Error())
				os.Exit(1)
			}
			exercises = track.Exercises
		}

		errs := []error{}
		for _, exercise := range exercises {
			readme, err := track.NewExerciseReadme(root, trackID, exercise.Slug)
			if err != nil {
				errs = append(errs, err)
				continue
			}
			if err := readme.Write(); err != nil {
				errs = append(errs, err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&genSlug, "only", "o", "", "Generate READMEs for just the exercise specified (by the slug).")
}
