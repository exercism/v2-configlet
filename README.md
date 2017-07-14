# Configlet

A linter for exercism language track repositories.

The API that delivers language-specific exercism problems makes
certain assumptions. Configlet makes it simple to verify up-front
that the changes to existing problems or the addition of new problems
will play nicely with the API.

Configlet checks for the following configuration issues:

1. `config.json` contents that are invalid according to [the specification](https://github.com/exercism/problem-specifications/blob/master/CONTRIBUTING.md#track-configuration-file).
1. Inconsistencies between the lists of track slugs in `config.json` and the corresponding implementation files:
    * Slugs referenced in `config.json` that have no corresponding implementation.
    * Slugs referenced in `config.json` whose implementation is missing an example solution.
    * Implementations for slugs that are not referenced in `config.json`.
    * Implementations for slugs that have been declared as foregone in `config.json`.

## Usage

```bash
$ configlet path/to/problem/repository
```

If you have [installed the configlet binary](https://github.com/exercism/configlet/releases/latest)
and are at the root of an exercism language track, then you can run the following:

```bash
$ configlet lint .
```
