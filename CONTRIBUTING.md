# Contributing Guide

First, thank you! :tada:
Exercism would be impossible without people like you being willing to spend time and effort making things better.

## Dependencies

You'll need Go version 1.7 or higher. Follow the directions on http://golang.org/doc/install

## Development

If you've never contributed to a Go project before this is going to feel a little bit foreign.

The TL;DR is: **don't clone your fork**, and it matters where on your filesystem the project gets cloned to.

If you don't care how and why and just want something that works, follow these steps:

1. [fork this repo][fork]
1. `go get github.com/exercism/configlet`
1. `cd $GOPATH/src/github.com/exercism/configlet`
1. `git remote set-url origin https://github.com/<your-github-username>/configlet`
1. `go get -u github.com/golang/dep/cmd/dep`
1. `dep ensure`

Then make the change as usual, and submit a pull request. Please provide tests for the changes where possible.

If you care about the details, check out the blog post [Contributing to Open Source Repositories in Go][contrib-blog] on the Splice blog.

## Running the Tests

To run the tests locally on Linux or MacOS, use

```
go test $(go list ./... | grep -v vendor)
```

On Windows, the command is more painful (sorry!):

```
for /f "" %G in ('go list ./... ^| find /i /v "/vendor/"') do @go test %G
```

As of Go 1.9 this is simplified to `go test ./...`.

### Cutting a release

This process could probably be somewhat automated, but for now it's manual.

1. Bump the version on master, commit with message "Bump version to vX.Y.Z" (actual version, though), and push to GitHub.
1. Run `bin/build` to cross-compile for all platforms. The binaries will be built into the `release` directory.
1. [Draft a new release](https://github.com/exercism/configlet/releases/new)
  * Select "recent commits" from the "Target" dropdown, then select the commit where you bumped the version.
  * Drag the releases from your release directory into the drop target in the form.
  * Look at the compare view between the previous release tag and the current master:
    https://github.com/exercism/configlet/compare/vX.Y.Z...master
  * Add a title that reflects the most important change.
  * Add a body that adds whatever detail seems relevant.
1. Click "publish release"

Travis will fetch the latest release automatically the next time it tries to build a track repository.

[fork]: https://github.com/exercism/configlet/fork
[contrib-blog]: https://splice.com/blog/contributing-open-source-git-repositories-go/
