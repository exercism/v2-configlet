package main

import "testing"

const fakeTrackPath = "./fixtures/track"

func TestTrackDirs(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	dirs, err := track.Dirs()
	if err != nil {
		t.Errorf("Borked. Don't understand.")
	}

	expected := []string{"amethyst", "beryl", "garnet", "ignored"}

	if len(dirs) != len(expected) {
		t.Errorf("Expected len(dirs)==%v to equal len(expected)==%v", dirs, expected)
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
	if err != nil {
		t.Errorf("Borked. Don't understand.")
	}

	expected := []string{"amethyst", "beryl", "crystal"}

	if len(problems) != len(expected) {
		t.Errorf("Expected len(problems)==%v to equal len(expected)==%v", problems, expected)
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
	if err != nil {
		t.Errorf("Borked. Don't understand.")
	}

	expected := []string{"amethyst", "bin", "beryl", "crystal", "ignored", "no-such-dir", ".git", "opal"}

	if len(slugs) != len(expected) {
		t.Errorf("Expected len(slugs)==%v to equal len(expected)==%v", slugs, expected)
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
	if err != nil {
		t.Errorf("Blew up: %v", err)
	}

	if len(problems) != 1 {
		t.Errorf("Expected len(%v)==1", problems)
	}

	if problems[0] != "crystal" {
		t.Errorf("Expected missing problem to be 'crystal', but was %s", problems[0])
	}
}

func TestProblemIsUnconfigured(t *testing.T) {
	track := NewTrack(fakeTrackPath)

	problems, err := track.UnconfiguredProblems()
	if err != nil {
		t.Errorf("Blew up: %v", err)
	}

	if len(problems) != 1 {
		t.Errorf("Expected len(%v)==1", problems)
	}

	if problems[0] != "garnet" {
		t.Errorf("Expected unconfigured problem to be 'garnet', but was %s", problems[0])
	}
}
