package ui

import "os"

func ExamplePrint() {
	Out = os.Stdout
	Print("this is a message")
	// Output:
	// -> this is a message
}

func ExamplePrintError() {
	ErrOut = os.Stdout
	PrintError("this is an error")
	// Output:
	// -> this is an error
}
