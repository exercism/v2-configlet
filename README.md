# Configlet

`configlet` is an internal command-line tool for managing Exercism language track repositories. The `confliglet` binary is [fetched from the track repos](https://github.com/exercism/request-new-language-track/blob/master/bin/fetch-configlet) and the linter is [executed during the build / CI process](https://github.com/exercism/request-new-language-track/blob/master/.travis.yml). In addition to linting, configlet provides a number of commands to assist in track creation and maintenance:

 * [Lint](#lint)
 * [Format](#format)
 * [Generate](#generate)
 * [Tree](#tree)
 * [Upgrade](#upgrade)
 * [UUID](#uuid)


## Install / Build

Configlet is a standalone Go application which compiles to a binary. Download
[the latest build](https://github.com/exercism/configlet/releases/latest), extract, and execute.

To build the application from source, run: `go build ./...`

To execute the tests, run: `go test ./...`


## Usage

```bash
configlet [command] <path/to/track>
```

If you are at the root of an exercism language track, then you can run the following:

```bash
configlet [command] .
```

## Lint

Exercism makes certain assumptions about language tracks. The configlet `lint` command makes it simple to verify up-front that the changes to a track's configuration, as well as changes and additions to the exercises will play nicely with the website.

`configlet lint` checks for the following configuration issues:

1. `config.json` contents that are invalid according to [the specification](https://github.com/exercism/problem-specifications/blob/master/CONTRIBUTING.md#track-configuration-file).
1. Inconsistencies between the lists of track slugs in `config.json` and the corresponding implementation files:
    * Slugs referenced in `config.json` that have no corresponding implementation.
    * Slugs referenced in `config.json` whose implementation is missing an example solution.
    * Implementations for slugs that are not referenced in `config.json`.
    * Implementations for slugs that have been declared as foregone in `config.json`.


## Format

Inspired by Go's [`gofmt`](https://blog.golang.org/go-fmt-your-code) tool, configlet's `fmt` command will consistently format a track's configuration files.

`configlet fmt` formats according to the following rules:

1. The JSON files, `config.json` and `maintainers.json` will be indented by 2 spaces.
1. In the `config.json` file:
    * Exercises will have their list of topics sorted alphabetically.
    * Topics names will be normalised to be lowercase and contain underscores in place of spaces.


## Generate

The configlet `generate` command may be used to generate README's for all exercises on the track, as a whole or individually. The overall purpose of this command is to utilize common [`problem-specifications`](https://github.com/exercism/problem-specifications) data in generating a *template-based* README for each exercise, while at the same time permitting overrides of this data on a track exercise basis.

Using this command for README generation allows for some conveniences:

1. The use of a track specific README template file.
1. The automatic inclusion of information from the `problem-specifications` repository.
1. Allow overrides of that `problem-specification` data on a per-exercise basis.
1. The ability to insert language specific exercise information ("hints") to an exercise README.
1. Full customization of a README template per exercise (if necessary).

`generate` looks for specific files in the track root's `config/` directory and in each exercise's `.meta/` directory. It will then use these files, if present, and the `problem-specifications` data to generate an exercise's README.

Much of the utility of this command comes from the ability to *locally override* README templates and exercise information.

(When working with READMEs you may find [a local renderer for GitHub Markdown](https://github.com/joeyespo/grip) helpful to preview your work before committing.)

### The README Template

The template file used as the basis for README generation lives in a track root's `config/` directory as [`config/exercise_readme.go.tmpl`](https://github.com/exercism/lua/blob/master/config/exercise_readme.go.tmpl). This template file may be overridden for an exercise by placing the overriding template in the exercises `.meta/readme.go.tmpl`.

As configlet is written in the [Go language](https://golang.org/) this README template file is in the [Go template format](https://golang.org/pkg/text/template/).

There are a number of template variable substitutions you may place in the template file:

#### .Spec.Description

This variable is sourced from an exercise's [`description.md`](https://github.com/exercism/problem-specifications/blob/master/exercises/hamming/description.md) file in the `problem-specifications` repo. You may override this variable's contents for an exercise by adding a `.meta/description.md` file in that track exercise's directory.

#### .Spec.Credits

The credits are a description of the source of an exercise with an optional hyperlink to that source. This information originates from the [`metadata.yml`](https://github.com/exercism/problem-specifications/blob/master/exercises/hamming/metadata.yml) located in the exercise's `problem-specifications` entry. You may override this information for an exercise by adding a `.meta/metadata.yml` file in that track exercise's directory.

#### .Spec.Name

This variable is a readable version of the exercise's slug in [title case](https://golang.org/pkg/strings/#Title). There are alternative formats of the name available. These formats may be useful if you need to [reference the exercise name](https://github.com/exercism/groovy/blob/1ffee8ea0df4492b349e367ac9ba88f1124bc038/config/exercise_readme.go.tmpl#L13) in regards to tooling.

| Variable            | Contents
| --------            | --------
| .Spec.Slug          | difference-of-squares
| .Spec.Name          | Difference Of Squares
| .Spec.MixedCaseName | DifferenceOfSquares
| .Spec.SnakeCaseName | difference\_of\_squares

#### .TrackInsert

Language tracks will most likely have some unique information common to every exercise in the track (testing, environment configuration, etc...). This may be placed in a track's `config/exercise-readme-insert.md` file the contents of which will then be available in this template variable.

#### .Hints

Exercises may have information specific to that exercise's implementation in the track language (for example, the introduction of a specific language concept). In this case placing a [`.meta/hints.md`](https://github.com/exercism/go/blob/nextercism/exercises/leap/.meta/hints.md) in that track exercise's directory will make those contents available in this template variable.


## Tree

The track configuration file can be hard to review. The `tree` command can help with the process of setting up your configuration file. It will:

1. Display the core track exercises and unlocked exercises as a tree.
1. List out the bonus exercises separately.
1. Issue warnings if expected elements from the configuration are missing.
1. Optionally show the difficulty of the exercises via the `--with-difficulty` option.


## Upgrade

The configlet `upgrade` command downloads and installs the latest released version of configlet. Running the upgrade command on an already up-to-date version of configlet will exit with no change to the system. The version command `configlet version -l` can be used to check for the latest available version.


## UUID

Exercises in each track config.json file must have a [UUID](https://en.wikipedia.org/wiki/Universally_unique_identifier). You may request a randomly generated exercise UUID from configlet using:

```bash
$ configlet uuid
78aa565f-632d-47c0-a190-5144c91d0e33
```
