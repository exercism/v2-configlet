package main

import (
	"fmt"
	"os"
)

// Check identifies configuration problems.
type Check func() ([]string, error)

// ConfigError defines the error message for a Check.
type ConfigError struct {
	check Check
	msg   string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: configlet path/to/problems/repository")
		os.Exit(1)
	}

	path := os.Args[1]
	fmt.Printf("Evaluating %s\n", path)

	hasErrors := Evaluate(path)

	if hasErrors {
		os.Exit(1)
	}

	fmt.Println("... OK")
}
