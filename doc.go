/*
Command configlet verifies Exercism's track configuration.

Each track has a script that downloads the binary to run on Travis CI.

See http://exercism.io for more information about the project.

Compile the CLI:

    $ go build -o out/configlet

Run the lint command:

    $ path/to/configlet lint path/to/track

Cross-compile using the build script:

    $ bin/build

To test the tool, use the fixture files:

    $ go build -o out/configlet
    $ out/configlet fixtures/numbers

The command should exit with status 1, and output a list of issues.
*/
package main
