package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: configlet path/to/problems/repository")
		os.Exit(1)
	}

	path := os.Args[1]

	fmt.Printf("Evaluating %s\n", path)

	track := NewTrack(path)

	hasErrors := false
	if !track.hasValidConfig() {
		hasErrors = true
		fmt.Println("-> config.json is invalid")
	}

	problems, err := track.MissingProblems()
	if err != nil {
		hasErrors = true
		fmt.Errorf("-> %v", err)
	}

	if len(problems) > 0 {
		hasErrors = true
		fmt.Printf("-> No directory found for %v.\n", problems)
	}

	problems, err = track.UnconfiguredProblems()
	if err != nil {
		hasErrors = true
		fmt.Errorf("-> %v", err)
	}

	if len(problems) > 0 {
		hasErrors = true
		fmt.Printf("-> config.json does not include %v.\n", problems)
	}

	problems, err = track.ProblemsLackingExample()
	if err != nil {
		hasErrors = true
		fmt.Errorf("-> %v", err)
	}

	if len(problems) > 0 {
		hasErrors = true
		fmt.Printf("-> missing example solution in %v.\n", problems)
	}

	problems, err = track.ForegoneViolations()
	if err != nil {
		hasErrors = true
		fmt.Print("-> %v", err)
	}

	if len(problems) > 0 {
		fmt.Printf("-> %v should not be implemented.\n", problems)
	}

	if hasErrors {
		os.Exit(1)
	}

	fmt.Println("... OK")
}
