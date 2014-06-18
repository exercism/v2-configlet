package configlet

import "testing"

func TestBrokenConfig(t *testing.T) {
	_, err := Load("./fixtures/broken.json")
	if err == nil {
		t.Errorf("Expected Load() to complain that it couldn't parse the JSON")
	}
}

func TestValidConfig(t *testing.T) {
	path := "./fixtures/valid.json"
	_, err := Load(path)
	if err != nil {
		t.Errorf("Config at %s should be valid, but barfed: %s", path, err)
	}
}
