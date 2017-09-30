package cmd

import (
	"os"
	"path/filepath"

	"github.com/exercism/configlet/ui"
)

func ExampleVisualize() {
	visualizeTrack(filepath.FromSlash("../fixtures/visualize/config.json"))
	// Output:
	// Numbers
	// =======
	//
	// core
	// ----
	// ├─ one
	// │  └─ five
	// │
	// ├─ two
	// │  ├─ six
	// │  │  ├─ ten
	// │  │  └─ twelve
	// │  └─ thirteen
	// │
	// ├─ nine
	// │
	// └─ eleven
	//
	// bonus
	// -----
	// seven
	// eight
}

func ExampleVisualizeDifficulty() {

	orig := showDifficulty
	showDifficulty = true
	defer func() { showDifficulty = orig }()

	visualizeTrack(filepath.FromSlash("../fixtures/visualize/config.json"))
	// Output:
	// Numbers
	// =======
	//
	// core
	// ----
	// ├─ one [1]
	// │  └─ five [2]
	// │
	// ├─ two [1]
	// │  ├─ six [3]
	// │  │  ├─ ten [6]
	// │  │  └─ twelve [8]
	// │  └─ thirteen [5]
	// │
	// ├─ nine [4]
	// │
	// └─ eleven [4]
	//
	// bonus
	// -----
	// seven [3]
	// eight [5]
}

func ExampleVisualizeWarnings() {
	orig := ui.ErrOut
	ui.ErrOut = os.Stdout

	defer func() { ui.ErrOut = orig }()

	visualizeTrack(filepath.FromSlash("../fixtures/visualize/config-outdated.json"))
	// Output:
	// Numbers
	// =======
	// -> Cannot find any unlockable exercises, this track may be missing a nextercism compatible configuration.
	// -> Cannot find any core exercises, this track may be missing a nextercism compatible configuration.
	//
	// bonus
	// -----
	// one
	// two
}
