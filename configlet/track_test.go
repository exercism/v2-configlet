package configlet

import (
	"sort"
	"strings"
	"testing"
)

const fakeTrackPath = "./fixtures/track"

func TestNewTrack(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	expected := []string{
		"fixtures/track/.git",
		"fixtures/track/beryl",
		"fixtures/track/bin",
		"fixtures/track/diamond",
		"fixtures/track/exercises/amethyst",
		"fixtures/track/exercises/melanite",
		"fixtures/track/ignored",
		"fixtures/track/sapphire",
	}

	var paths []string
	for _, path := range track.dirs {
		paths = append(paths, path)
	}
	if len(paths) != len(expected) {
		t.Errorf("len(paths)=%d - expected %d", len(paths), len(expected))
	}

	sort.Strings(paths)
	sort.Strings(expected)

	for i := 0; i < len(paths); i++ {
		if paths[i] != expected[i] {
			t.Errorf("got: %s, want: %s", paths[i], expected[i])
		}
	}
}

func TestTrackDirs(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	dirs, err := track.Dirs()
	assertNoError(t, err)

	expected := []string{
		".git",
		"amethyst",
		"beryl",
		"bin",
		"garnet",
		"ignored",
		"diamond",
		"melanite",
		"sapphire",
	}

	if len(dirs) != len(expected) {
		msg := "Expected len(dirs:%d)==%v to equal len(expected:%v)==%d"
		t.Errorf(msg, len(dirs), dirs, expected, len(expected))
	}

	for _, gemstone := range expected {
		_, ok := dirs[gemstone]
		if !ok {
			t.Errorf("Expected track.Dirs() to contain %s", gemstone)
		}
	}
}

func TestTrackProblems(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	problems, err := track.Problems()
	assertNoError(t, err)

	expected := []string{
		"amethyst",
		"beryl",
		"crystal",
		"melanite",
		"sapphire",
	}

	if len(problems) != len(expected) {
		msg := "Expected len(problems:%v)==%d to equal len(expected:%v)==%d"
		t.Errorf(msg, problems, len(problems), expected, len(expected))
	}

	for _, gemstone := range expected {
		_, ok := problems[gemstone]
		if !ok {
			t.Errorf("Expected track.Problems() to contain %s", gemstone)
		}
	}
}

func TestSlugs(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	slugs, err := track.Slugs()
	assertNoError(t, err)

	expected := []string{
		".git",
		"amethyst",
		"beryl",
		"bin",
		"crystal",
		"diamond",
		"ignored",
		"melanite",
		"no-such-dir",
		"opal",
		"pearl",
		"sapphire",
	}

	if len(slugs) != len(expected) {
		msg := "Expected len(slugs:%v)==%d to equal len(expected:%v)==%d"
		t.Errorf(msg, slugs, len(slugs), expected, len(expected))
	}

	for _, slug := range expected {
		_, ok := slugs[slug]
		if !ok {
			t.Errorf("Expected track.Slugs() to contain %s", slug)
		}
	}
}

func TestProblemIsMissing(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	problems, err := track.MissingProblems()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, problems, len(problems))
	}

	if problems[0] != "crystal" {
		t.Errorf("Expected missing problem to be 'crystal', but was %s", problems[0])
	}
}

func TestProblemIsUnconfigured(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	problems, err := track.UnconfiguredProblems()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, problems, len(problems))
	}

	if problems[0] != "garnet" {
		t.Errorf("Expected unconfigured problem to be 'garnet', but was %s", problems[0])
	}
}

func TestProblemLacksExample(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	problems, err := track.ProblemsLackingExample()
	assertNoError(t, err)

	if len(problems) != 2 {
		msg := "Expected len(problems)==2, but len(%v)==%d"
		t.Fatalf(msg, problems, len(problems))
	}

	sort.Strings(problems)

	if problems[0] != "beryl" {
		t.Errorf("Expected missing example to be on 'beryl' problem, but was %s", problems[0])
	}
	if problems[1] != "melanite" {
		t.Errorf("Expected missing example to be on 'melanite' problem, but was %s", problems[1])
	}
}

func TestForegoneViolations(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	problems, err := track.ForegoneViolations()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, problems, len(problems))
	}

	if problems[0] != "diamond" {
		t.Errorf("Expected violation to be 'diamond', but was %s", problems[0])
	}
}

func TestDuplicateSlugs(t *testing.T) {
	track, err := NewTrack(fakeTrackPath)
	assertNoError(t, err)

	problems, err := track.DuplicateSlugs()
	assertNoError(t, err)

	if len(problems) != 3 {
		msg := "Expected len(problems)==3, but len(%v)==%d"
		t.Errorf(msg, len(problems), problems)
	}

	expected := []string{"amethyst", "beryl", "crystal"}
	sort.Strings(problems)

	if strings.Join(problems, " ") != strings.Join(expected, " ") {
		t.Errorf("Expected duplicates to be '[amethyst beryl crystal]', but was %v", problems)
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("FAIL: %v", err)
	}
}
