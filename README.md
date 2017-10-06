# Configlet

A tool for managing Exercism language track repositories.

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
    * Topics names will be normalised to be lowercase and contain hyphens in place of spaces.

## Tree

The track configuration file can be hard to review, especially the new structure being used for [nextercism](https://github.com/exercism/prototype). The `tree` command can help with the process of setting up your configuration file for nextercism. It will:

1. Display the core track exercises and unlocked exercises as a tree.
1. List out the bonus exercises separately.
1. Issue warnings if expected elements from the nextercism-style configuration are missing.
1. Optionally show the difficulty of the exercises via the `--difficulty` option.

### Usage

```bash
$ configlet [command] <path/to/track>
```

If you have [installed the configlet binary](https://github.com/exercism/configlet/releases/latest)
and are at the root of an exercism language track, then you can run the following:

```bash
$ configlet [command] .
```

