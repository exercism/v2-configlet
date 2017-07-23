package track

import (
	"io/ioutil"
	"path/filepath"
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

	for _, file := range files {
		if file.IsDir() {
			fp := filepath.Join(dir, file.Name())

			ex, err := NewExercise(fp, track.Config.PatternGroup)
			if err != nil {
				return track, err
			}

			track.Exercises = append(track.Exercises, ex)
		}
	}
	return track, nil
}
