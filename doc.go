/*
Command configlet sanity-checks Exercism track configuration.

Each track has a script that downloads the binary to run on Travis CI.

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

The command should exit with status 1, and output a list of issues.

*/
package main
