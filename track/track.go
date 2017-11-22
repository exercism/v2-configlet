package track

import (
	"io/ioutil"
	"path/filepath"
)

// Track is a collection of Exercism exercises for a programming language.
type Track struct {
	ID               string
	Config           Config
	MaintainerConfig MaintainerConfig
	Exercises        []Exercise
	path             string
}

// ExerciseSlugs returns all implemented exercise slugs of track
func (t Track) ExerciseSlugs() []string {
	var slugs []string
	for _, e := range t.Exercises {
		slugs = append(slugs, e.Slug)
	}
	return slugs
}

// New loads a track.
func New(path string) (Track, error) {
	track := Track{
		path: filepath.FromSlash(path),
	}

	ap, err := filepath.Abs(track.path)
	if err != nil {
		return track, err
	}
	track.ID = filepath.Base(ap)

	c, err := NewConfig(filepath.Join(path, "config.json"))
	if err != nil {
		return track, err
	}
	track.Config = c

	mc, err := NewMaintainerConfig(filepath.Join(path, "config", "maintainers.json"))
	if err != nil {
		return track, err
	}
	track.MaintainerConfig = mc

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
