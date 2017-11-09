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
	Slug          string
	SolutionPath  string
	TestSuitePath string
}

// Exercises is a slice of Exercise to define functions
type Exercises []Exercise

// Fold iterates over slice, runs given function and collects slugs
func (es Exercises) Fold(isValid func(Exercise) bool) (valid []string, invalid []string) {
	for _, e := range es {
		if isValid(e) {
			valid = append(valid, e.Slug)
		} else {
			invalid = append(invalid, e.Slug)
		}
	}
	return
}

// NewExercise loads an exercise.
func NewExercise(root string, pg PatternGroup) (Exercise, error) {
	ex := Exercise{
		Slug: filepath.Base(root),
	}

	err := setPath(root, pg.SolutionPattern, &ex.SolutionPath)
	if err != nil {
		return ex, err
	}

	err = setPath(root, pg.TestPattern, &ex.TestSuitePath)
	if err != nil {
		return ex, err
	}

	return ex, err
}

// setPath sets the value of field to the file path matched by pattern.
// The resulting file path, if matched, will be relative to root.
func setPath(root, pattern string, field *string) error {

	if pattern == "" {
		return nil
	}

	rgx, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if rgx.Match([]byte(path)) {
			prefix := fmt.Sprintf("%s%s", root, string(filepath.Separator))
			*field = strings.Replace(path, prefix, "", 1)
		}
		return nil
	}

	return filepath.Walk(root, walkFn)
}

// HasTestSuite checks that an exercise has a test suite.
func (ex Exercise) HasTestSuite() bool {
	return ex.TestSuitePath != ""
}

// IsValid checks that an exercise has a sample solution.
func (ex Exercise) IsValid() bool {
	return ex.SolutionPath != ""
}
