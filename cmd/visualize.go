package cmd

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"

	"github.com/exercism/configlet/track"
	"github.com/exercism/configlet/ui"
)

// Characters and spacing for the tree structure output:
const (
	indent     = 2   // how much to indent depth in the tree
	trunk      = "│" // descent in the tree
	branch     = "─" // prefix for an exercise
	fork       = "├" // normal trunk where a branch then starts (an exercise is listed)
	terminator = "└" // special fork where there is nothing below it
)

// configurationWarning is boilerplate that will be appended to any warning
// regarding potential mis-configuration for nextercism.
const configurationWarning = ", this track may be missing a nextercism compatible configuration."

// configPathExample is used in the Cobra usage output where the configfile path
// would be.
const configPathExample = "<path/to/config.json>"

// treeSpacing and treeBranching will be user over and over again
// in the tree creation, so create them once here.
var treeSpacing = strings.Repeat(" ", indent)

var treeBranching = strings.Repeat(branch, indent-1)

// showDifficulty holds --difficulty flag value indicates that we
// should display exercise difficulty after slug, by default we do not.
var showDifficulty bool

// visCmd defines the visualize command
var visCmd = &cobra.Command{
	Use:   "visualize " + configPathExample,
	Short: "View the track structure as a tree",
	Long: `The visualize command displays the track in a tree format, with core
exercises at root and unlocks located under their locking exercises. 

Bonus exercises are left in a list at the bottom after the tree display.

example output:

Go
--

├── hello-world
│
├── hamming
│   ├── nucleotide-count
│   └── rna-transcription
...

`,
	Example: fmt.Sprintf("  %s visualize %s --difficulty", binaryName, configPathExample),
	Run:     runVisualizer,
	Args:    cobra.ExactArgs(1),
}

// slugToExercise is a global lookup table of slugs to exercises
var slugToExercise = map[string]*exerciseUnlocks{}

// exerciseUnlocks is an extension of the exercise metatata type
// with with the exercises it unlocks
type exerciseUnlocks struct {
	track.ExerciseMetadata
	Unlocks []string // slugs of unlocked exercises
}

// getDescription is a utility that will return the description for
// an exerciseUnlock with the difficulty appended if the --difficulty
// flag was set
func (e exerciseUnlocks) getDescription() string {

	if showDifficulty {
		return fmt.Sprintf("%s [%d]", e.Slug, e.Difficulty)
	}

	return e.Slug
}

// writeln is a convenience wrapper around printing the line directly
// in case we (eventually) need to do something else with the string
// or writer beforehand.
func writeln(s string) {
	fmt.Fprintln(os.Stdout, s)
}

// writelns will write an array of lines using writeln.
func writelns(ss []string) {
	for _, s := range ss {
		writeln(s)
	}
}

// runVisualizer kicks off the visualization and will print any
// errors from the process.
func runVisualizer(cmd *cobra.Command, args []string) {
	for _, arg := range args {
		if err := visualizeTrack(arg); err != nil {
			ui.PrintError(err)
		}
	}
}

// printConfigurationWarning is a utility that will print s to Error,
// but will prefix it with a nextercism specific configuration warning
// use it if you expect something via nextercism but it is not present in
// the configuration.
func printConfigurationWarning(s string) {
	ui.PrintError(s + configurationWarning)
}

// tree is responsible for outputting the actual tree structure
// it will operate recursively on the unlocks and send contents directly
// to output.
//
// isLast is a special indicator, callers can use this to note that
// the exercise being processed is the last in a sequence, this will
// make some special tweaks to the output format to look a little
// more pleasant.
func tree(e *exerciseUnlocks, depth int, isLast bool) {
	var buffer bytes.Buffer // Holds for the generated output of this exercise.

	numChildren := len(e.Unlocks) // Unlocks are children in the tree context.
	hasChildren := numChildren > 0

	// Create the pre-fixing for this exercise using depth to move
	// this exercise further and further to the right.
	for i := 0; i < depth; i++ {
		buffer.WriteString(trunk) // continue trunk from the parent exercise(s)
		buffer.WriteString(treeSpacing)
	}

	// Normally show the fork indicating a peer below unless there
	// is none.
	if !hasChildren && isLast {
		buffer.WriteString(terminator)
	} else {
		buffer.WriteString(fork)
	}

	// Show the exercise name, it will have the standard branch prefix.
	buffer.WriteString(treeBranching)
	buffer.WriteString(" ")
	buffer.WriteString(e.getDescription())
	writeln(buffer.String())

	// Now go into the children unlocks and do this all over again.
	for i, slug := range e.Unlocks {
		child := slugToExercise[slug]
		tree(child, depth+1, i == (numChildren-1))
	}

	// If depth is 0 (we are at root of the tree) we will add a little
	// extra spacing between this and the next exercise...
	// ...except for the last element because there is nothing below it
	// to space out.
	if depth == 0 && !isLast {
		writeln(trunk)
	}
}

func visualizeTrack(path string) error {

	// exercises is a list of all non-deprecated exercises, in config order.
	exercises := make([]exerciseUnlocks, 0)

	// coreExercises are slugs of exercises in core, in config order.
	coreExercises := make([]string, 0)

	// bonusExercises are slugs of exercises determined to be bonuses:
	// non-core, no unlocks, in config order.
	bonusExercises := make([]string, 0)

	config, err := track.NewConfig(path)
	if err != nil {
		return err
	}

	// Print a header: the language name with markdown style h1 underlining.
	writeln(config.Language)
	writeln(strings.Repeat("=", utf8.RuneCountInString(config.Language)))

	// Initial scan through the exercises for this track: filter out deprecated
	// exercises and setup the data structures from above.
	for _, e := range config.Exercises {
		// Completely ignore deprecated exercises. They are dead to us.
		if e.IsDeprecated {
			continue
		}
		// Container for this exercise
		eu := exerciseUnlocks{
			e,
			make([]string, 0), // Our unlock slugs, filled in on second pass.
		}

		exercises = append(exercises, eu)

		// Add to slug based global lookup table
		slugToExercise[e.Slug] = &eu

		if eu.IsCore {
			coreExercises = append(coreExercises, eu.Slug)
		} else if eu.UnlockedBy == "" {
			bonusExercises = append(bonusExercises, eu.Slug)
		}
	}

	// Second pass through exercises, fill out the unlocks.
	// Look to see if there are no unlocks, (unlocksPresent is never set to
	// true) if so issue a nextercism warning.
	unlocksPresent := false
	for _, e := range exercises {
		if e.UnlockedBy == "" {
			continue
		}

		unlocksPresent = true
		parent := slugToExercise[e.UnlockedBy]
		parent.Unlocks = append(parent.Unlocks, e.Slug)
	}

	// This is more of a warning than an error. No stdwarn :(
	if !unlocksPresent {
		printConfigurationWarning("Cannot find any unlockable exercises")
	}

	// If we have core exercises add markdown-style secondary header then
	// loop through and show in tree format.
	// If we don't have any core exercises, warn about configuration.
	numCore := len(coreExercises)
	if numCore > 0 {
		writelns([]string{"", "core", "----"})

		lastSlug := coreExercises[numCore-1] // Used to set the isLast hint.
		for _, slug := range coreExercises {
			e := slugToExercise[slug]
			tree(e, 0, (slug == lastSlug))
		}
	} else {
		printConfigurationWarning("Cannot find any core exercises")
	}

	// This is not a tree structure, so unlike core above we do not use the
	// tree output. Just a normal listing. Otherwise this is like the core loop.
	numBonus := len(bonusExercises)
	if numBonus > 0 {
		writelns([]string{"", "bonus", "-----"})

		for _, slug := range bonusExercises {
			e := slugToExercise[slug]
			writeln(e.getDescription())
		}
	} else {
		printConfigurationWarning("Cannot find any bonus exercises")
	}

	return nil
}

func init() {
	RootCmd.AddCommand(visCmd)
	visCmd.Flags().BoolVar(&showDifficulty, "difficulty", false, "display the difficulty of the exercises")
}
