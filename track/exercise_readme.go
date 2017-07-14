package track

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
)

const (
	dirExercises   = "exercises"
	filenameReadme = "README.md"
)

var (
	pathTrackTemplate            = filepath.Join("config", "exercise_readme.go.tmpl")
	pathTrackInsert              = filepath.Join("config", "exercise-readme-insert.md")
	pathTrackInsertDeprecated    = filepath.Join("docs", "EXERCISE_README_INSERT.md")
	pathExerciseTemplate         = filepath.Join(".meta", "readme.go.tmpl")
	pathExerciseInsert           = filepath.Join(".meta", "hints.md")
	pathExerciseInsertDeprecated = "HINTS.md"
)

type ExerciseReadme struct {
	Spec        *ProblemSpecification
	Hints       string
	TrackInsert string
	template    string
	trackDir    string
	dir         string
}

func NewExerciseReadme(root, trackID, slug string) (ExerciseReadme, error) {
	readme := ExerciseReadme{
		trackDir: filepath.Join(root, trackID),
		dir:      filepath.Join(root, trackID, dirExercises, slug),
	}

	if err := readme.readTemplate(); err != nil {
		return readme, err
	}

	spec, err := NewProblemSpecification(root, trackID, slug)
	if err != nil {
		spec = &ProblemSpecification{}
	}
	readme.Spec = spec

	if err := readme.readHints(); err != nil {
		return readme, err
	}

	if err := readme.readTrackInsert(); err != nil {
		return readme, err
	}

	return readme, nil
}

func (readme ExerciseReadme) Generate() (string, error) {
	t, err := template.New("readme").Parse(readme.template)
	if err != nil {
		return "", err
	}

	var bb bytes.Buffer
	t.Execute(&bb, readme)
	return bb.String(), nil
}

func (readme ExerciseReadme) Write() error {
	s, err := readme.Generate()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(readme.dir, filenameReadme), []byte(s), 0644)
}

func (readme *ExerciseReadme) readTrackInsert() error {
	b, err := ioutil.ReadFile(filepath.Join(readme.trackDir, pathTrackInsert))
	if err == nil {
		readme.TrackInsert = string(b)
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	b, err = ioutil.ReadFile(filepath.Join(readme.trackDir, pathTrackInsertDeprecated))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	readme.TrackInsert = string(b)
	return nil
}

func (readme *ExerciseReadme) readHints() error {
	b, err := ioutil.ReadFile(filepath.Join(readme.dir, pathExerciseInsert))
	if err == nil {
		readme.Hints = string(b)
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}
	b, err = ioutil.ReadFile(filepath.Join(readme.dir, pathExerciseInsertDeprecated))
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	readme.Hints = string(b)
	return nil
}

func (readme *ExerciseReadme) readTemplate() error {
	b, err := ioutil.ReadFile(filepath.Join(readme.dir, pathExerciseTemplate))
	if err == nil {
		readme.template = string(b)
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	b, err = ioutil.ReadFile(filepath.Join(readme.trackDir, pathTrackTemplate))
	if err != nil {
		return err
	}
	readme.template = string(b)
	return nil
}
