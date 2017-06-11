package main

func ExampleEvaluate() {
	Evaluate("./configlet/fixtures/track")
	// Output:
	// -> An exercise with slug 'crystal' is referenced in config.json, but no implementation was found.
	// -> An implementation for 'garnet' was found, but config.json does not reference this exercise.
	// -> The implementation for 'beryl' is missing an example solution.
	// -> The implementation for 'melanite' is missing an example solution.
	// -> An implementation for 'diamond' was found, but config.json specifies that it should be foregone (not implemented).
	// -> The exercise 'amethyst' was found in multiple (conflicting) categories in config.json.
	// -> The exercise 'crystal' was found in multiple (conflicting) categories in config.json.
}
