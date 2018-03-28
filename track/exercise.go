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
	ReadmePath    string
	SolutionPath  string
	TestSuitePath string
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

	err = setPath(root, "README\\.md", &ex.ReadmePath)
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

// HasReadme checks that an exercise has a README.
func (ex Exercise) HasReadme() bool {
	return ex.ReadmePath != ""
}

// HasTestSuite checks that an exercise has a test suite.
func (ex Exercise) HasTestSuite() bool {
	return ex.TestSuitePath != ""
}

// IsValid checks that an exercise has a sample solution.
func (ex Exercise) IsValid() bool {
	return ex.SolutionPath != ""
}
