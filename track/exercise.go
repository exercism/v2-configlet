package track

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Exercise is an implementation of an Exercism exercise.
type Exercise struct {
	Slug         string
	SolutionPath string
}

// NewExercise loads an exercise.
func NewExercise(root string, rgx *regexp.Regexp) (Exercise, error) {
	ex := Exercise{
		Slug: filepath.Base(root),
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if rgx.Match([]byte(path)) {
			prefix := fmt.Sprintf("%s%s", root, string(filepath.Separator))
			ex.SolutionPath = strings.Replace(path, prefix, "", 1)
		}
		return nil
	}

	err := filepath.Walk(root, walkFn)
	return ex, err
}

// IsValid checks that an exercise has a sample solution.
func (ex Exercise) IsValid() bool {
	return ex.SolutionPath != ""
}
