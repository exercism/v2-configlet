package ui

import (
	"fmt"
	"io"
	"os"
)

var (
	// prefix is prepended to ui messages
	prefix = "->"
	// Out is a writer for messages
	Out io.Writer = os.Stdout
	// ErrOut is a writer for error messages
	ErrOut io.Writer = os.Stderr
)

// Print writes msg to Out
func Print(msg ...interface{}) {
	printer(Out, msg...)
}

// PrintError writes msg to ErrOut
func PrintError(msg ...interface{}) {
	printer(ErrOut, msg...)
}

func printer(w io.Writer, msg ...interface{}) {
	a := append([]interface{}{prefix}, msg...)
	fmt.Fprintln(w, a...)
}
