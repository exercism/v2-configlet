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

	track := NewTrack(os.Args[1])

	hasErrors := false
	if !track.hasValidConfig() {
		hasErrors = true
		fmt.Println("config.json is invalid")
	}

	problems, err := track.MissingProblems()
	if err != nil {
		hasErrors = true
		fmt.Errorf("%v", err)
	}

	if len(problems) > 0 {
		hasErrors = true
		fmt.Printf("No directory found for %v.\n", problems)
	}

	problems, err = track.UnconfiguredProblems()
	if err != nil {
		hasErrors = true
		fmt.Errorf("%v", err)
	}

	if len(problems) > 0 {
		hasErrors = true
		fmt.Printf("config.json does not include %v\n", problems)
	}

	if hasErrors {
		os.Exit(1)
	}

	fmt.Println("OK")
}
