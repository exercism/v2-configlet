package main

import (
	"flag"
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

const (
	// Version represents the latest released version of the project.
	// Configlet follows semantic versioning.
	Version = "2.0.0"
)

var showVersion = flag.Bool("version", false, "output the version of the tool")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("configlet v%s\n", Version)
		os.Exit(0)
	}

	if len(os.Args) != 2 {
		fmt.Println("Usage: configlet path/to/track/repository")
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
