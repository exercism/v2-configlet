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
	Version = "2.0.1"
)

var showVersion = flag.Bool("version", false, "output the version of the tool")

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("configlet v%s\n", Version)
		os.Exit(0)
	}

	var path string
	switch len(os.Args) {
	case 3:
		path = os.Args[2]
	case 2:
		path = os.Args[1]
	default:
		fmt.Println("Usage: configlet lint path/to/track")
		os.Exit(1)
	}

	fmt.Printf("Evaluating %s\n", path)

	hasErrors := Evaluate(path)

	if hasErrors {
		os.Exit(1)
	}

	fmt.Println("... OK")
}
