package configlet

import "testing"

const fakeTrackPath = "./fixtures/track"

func TestTrackDirs(t *testing.T) {
	track := NewTrack(fakeTrackPath)

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
	}

	if len(dirs) != len(expected) {
		msg := "Expected len(dirs:%d)==%v to equal len(expected:%d)==%v"
		t.Errorf(msg, len(dirs), dirs, len(expected), expected)
	}

	for _, gemstone := range expected {
		_, ok := dirs[gemstone]
		if !ok {
			t.Errorf("Expected track.Dirs() to contain %s", gemstone)
		}
	}
}

func TestTrackProblems(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	problems, err := track.Problems()
	assertNoError(t, err)

	expected := []string{
		"amethyst",
		"beryl",
		"crystal",
	}

	if len(problems) != len(expected) {
		msg := "Expected len(problems:%d)==%v to equal len(expected:%d)==%v"
		t.Errorf(msg, len(problems), problems, len(expected), expected)
	}

	for _, gemstone := range expected {
		_, ok := problems[gemstone]
		if !ok {
			t.Errorf("Expected track.Problems() to contain %s", gemstone)
		}
	}
}

func TestSlugs(t *testing.T) {
	track := NewTrack(fakeTrackPath)

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
		"no-such-dir",
		"opal",
		"pearl",
	}

	if len(slugs) != len(expected) {
		msg := "Expected len(slugs:%d)==%v to equal len(expected:%d)==%v"
		t.Errorf(msg, len(slugs), slugs, len(expected), expected)
	}

	for _, slug := range expected {
		_, ok := slugs[slug]
		if !ok {
			t.Errorf("Expected track.Slugs() to contain %s", slug)
		}
	}
}

func TestProblemIsMissing(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	problems, err := track.MissingProblems()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, len(problems), problems)
	}

	if problems[0] != "crystal" {
		t.Errorf("Expected missing problem to be 'crystal', but was %s", problems[0])
	}
}

func TestProblemIsUnconfigured(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	problems, err := track.UnconfiguredProblems()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, len(problems), problems)
	}

	if problems[0] != "garnet" {
		t.Errorf("Expected unconfigured problem to be 'garnet', but was %s", problems[0])
	}
}

func TestProblemLacksExample(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	problems, err := track.ProblemsLackingExample()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, len(problems), problems)
	}

	if problems[0] != "beryl" {
		t.Errorf("Expected missing example to be on 'beryl' problem, but was %s", problems[0])
	}
}

func TestForegoneViolations(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	problems, err := track.ForegoneViolations()
	assertNoError(t, err)

	if len(problems) != 1 {
		msg := "Expected len(problems)==1, but len(%v)==%d"
		t.Errorf(msg, len(problems), problems)
	}

	if problems[0] != "diamond" {
		t.Errorf("Expected violation to be 'diamond', but was %s", problems[0])
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("FAIL: %v", err)
	}
}
