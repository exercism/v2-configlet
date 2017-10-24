package cmd

import (
	"os"
	"path/filepath"

	"github.com/exercism/configlet/ui"
)

func ExampleTree() {
	treeTrack(filepath.FromSlash("../fixtures/tree/config.json"))
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

func ExampleTreeDifficulty() {

	orig := withDifficulty
	withDifficulty = true
	defer func() { withDifficulty = orig }()

	treeTrack(filepath.FromSlash("../fixtures/tree/config.json"))
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

func ExampleTreeWarnings() {
	orig := ui.ErrOut
	ui.ErrOut = os.Stdout

	defer func() { ui.ErrOut = orig }()

	treeTrack(filepath.FromSlash("../fixtures/tree/config-outdated.json"))
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

func ExampleInvalidUnlockedBy(){
	orig := ui.ErrOut
	ui.ErrOut = os.Stdout

	defer func() { ui.ErrOut = orig }()

	treeTrack(filepath.FromSlash("../fixtures/tree/config-invalid-unlocked-by.json"))
	// Output:
	// Numbers
	// =======
	// -> Exercise "six" has an invalid unlocked_by slug: "a non-existing exercise", this track may be missing a nextercism compatible configuration.
	//
	// core
	// ----
	// ├─ one
	// │  └─ five
	// │
	// └─ two
	//
	// bonus
	// -----
	// seven
}