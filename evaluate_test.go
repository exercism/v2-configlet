package main

func ExampleEvaluate() {
	Evaluate("./configlet/fixtures/track")
	// Output:
	// -> No directory found for [crystal].
	// -> config.json does not include [garnet].
	// -> missing example solution in [beryl].
	// -> [diamond] should not be implemented.
	// -> [amethyst beryl crystal] found in multiple categories.
}
