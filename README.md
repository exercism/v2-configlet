# Configlet

A linter for exercism language track repositories.

The API that delivers language-specific exercism problems makes
certain assumptions. Configlet makes it simple to verify up-front
that the changes to existing problems or the addition of new problems
will play nicely with the API.

There are three common problems that occur:

1. The `config.json` [(documented here)](https://github.com/exercism/x-common/blob/master/CONTRIBUTING.md#track-configuration-file) might be invalid.
1. Problems might be missing a reference solution.
1. Problems might be implemented (test suite + reference solution), but not configured.
1. Slugs in the configuration might not have a corresponding problem.

## Usage

```bash
$ configlet path/to/problem/repository
```

If you have [installed the configlet binary](https://github.com/exercism/configlet/releases/latest)
and are at the root of an exercism language track, then you can run the following:

```bash
$ configlet .
```
