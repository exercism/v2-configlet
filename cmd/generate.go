package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/exercism/configlet/track"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var (
	genSlug  string
	specPath string
)

var (
	// generateCmd represents the generate command
	generateCmd = &cobra.Command{
		Use:   "generate <path/to/track>",
		Short: "Generate exercise READMEs for an Exercism language track",
		Long:  `Generate READMEs for Exercism exercises based on the contents of various files.`,
		Example: `  configlet generate <path/to/track> --only hello-world

  configlet generate <path/to/track> --spec-path <path/to/problem-specifications>
`,
		Run:  generate,
		Args: cobra.MinimumNArgs(1),
	}
)

func generate(cmd *cobra.Command, args []string) {
	path, err := filepath.Abs(filepath.FromSlash(args[0]))
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	root := filepath.Dir(path)
	trackID := filepath.Base(path)

	track.ProblemSpecificationsPath = filepath.Join(root, track.ProblemSpecificationsDir)
	if specPath != "" {
		track.ProblemSpecificationsPath = specPath
	}

	if _, err := os.Stat(track.ProblemSpecificationsPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "path not found: %s\n", track.ProblemSpecificationsPath)
		os.Exit(1)
	}

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

	errs := &multierror.Error{}
	for _, exercise := range exercises {
		readme, err := track.NewExerciseReadme(root, trackID, exercise.Slug)
		if err != nil {
			errs = multierror.Append(errs, err)
			continue
		}

		if err := readme.Write(); err != nil {
			errs = multierror.Append(errs, err)
		}
	}

	if err := errs.ErrorOrNil(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

}

func init() {
	RootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVarP(&genSlug, "only", "o", "", "Generate READMEs for just the exercise specified (by the slug).")
	generateCmd.Flags().StringVarP(&specPath, "spec-path", "p", "", "The location of the problem-specifications directory.")
}
