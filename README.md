# Configlet

A linter for exercism language track repositories.

Exercism makes certain assumptions about language tracks. Configlet makes it simple to verify up-front that the changes to existing exercises, or the addition of new exercises will play nicely with the website.

Configlet checks for the following configuration issues:

1. `config.json` contents that are invalid according to [the specification](https://github.com/exercism/x-common/blob/master/CONTRIBUTING.md#track-configuration-file).
1. Inconsistencies between the lists of track slugs in `config.json` and the corresponding implementation files:
    * Slugs referenced in `config.json` that have no corresponding implementation.
    * Slugs referenced in `config.json` whose implementation is missing an example solution.
    * Implementations for slugs that are not referenced in `config.json`.
    * Implementations for slugs that have been declared as foregone in `config.json`.

## Usage

```bash
$ configlet lint path/to/track
```

If you have [installed the configlet binary](https://github.com/exercism/configlet/releases/latest)
and are at the root of an exercism language track, then you can run the following:

```bash
$ configlet lint .
```
