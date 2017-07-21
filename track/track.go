package track

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

// Track is a collection of Exercism exercises for a programming language.
type Track struct {
	path      string
	Config    Config
	Exercises []Exercise
}

// New loads a track.
func New(path string) (Track, error) {
	track := Track{
		path: filepath.FromSlash(path),
	}

	c, err := NewConfig(filepath.Join(path, "config.json"))
	if err != nil {
		return track, err
	}
	track.Config = c

	dir := filepath.Join(track.path, "exercises")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return track, err
	}

	rgx, err := regexp.Compile(track.Config.SolutionPattern)
	for _, file := range files {
		if file.IsDir() {
			ex, err := NewExercise(filepath.Join(dir, file.Name()), rgx)
			if err != nil {
				return track, err
			}
			track.Exercises = append(track.Exercises, ex)
		}
	}
	return track, nil
}
