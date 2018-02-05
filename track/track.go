package track

import (
	"io/ioutil"
	"path/filepath"
	"regexp"
)

// Track is a collection of Exercism exercises for a programming language.
type Track struct {
	ID               string
	Config           Config
	MaintainerConfig MaintainerConfig
	Exercises        []Exercise
	path             string
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

	// Valid exercise directory names do not begin with `.` or `_`.
	re := regexp.MustCompile("^[._]")
	for _, file := range files {
		if file.IsDir() {
			fn := file.Name()
			if re.MatchString(fn) {
				continue
			}
			fp := filepath.Join(dir, fn)

			ex, err := NewExercise(fp, track.Config.PatternGroup)
			if err != nil {
				return track, err
			}

			track.Exercises = append(track.Exercises, ex)
		}
	}
	return track, nil
}
