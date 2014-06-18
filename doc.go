/*
Command configlet sanity-checks Exercism track configuration.

The binary is distributed with each track repository. The repositories are located at
github.com/exercism/x<TRACK>.

The command also runs on Travis-CI.

See http://exercism.io for more information about the project.

Compile the CLI:

    $ go build -o out/configlet

Run the command:

    $ path/to/configlet path/to/track

Cross-compile using the build script:

    $ bin/build

To verify the command, use the fixture files:

    $ go build -o out/configlet
    $ out/configlet configlet/fixtures/track

The command should exit with status 1, and output:

    Evaluating configlet/fixtures/track
    -> No directory found for [crystal].
    -> config.json does not include [garnet].
    -> missing example solution in [beryl].
    -> [diamond] should not be implemented.

*/
package main
