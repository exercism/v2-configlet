# Cutting a Configlet Release

`configlet` uses [GoReleaser](https://goreleaser.com) to automate the
release process. 

## Requirements

1. [Install GoReleaser](https://goreleaser.com/install/)
1. [Setup GitHub token](https://goreleaser.com/environment/#github-token)
1. Have a gpg key installed on your machine - it is [used for signing the artifacts](https://goreleaser.com/sign/)

## Confirm / Update the Changelog

Make sure all the recent changes are reflected in the "next release" section of the CHANGELOG.md file.  All the changes in the "next release" section should be moved to a new section that describes the version number, and gives it a date.

You can view changes using the /compare/ view:
https://github.com/exercism/configlet/compare/$PREVIOUS_RELEASE...master

Some features to consider - GoReleaser supports the [auto generation of a changelog](https://goreleaser.com/customization/#customize-the-changelog) and has a [release notes feature](https://goreleaser.com/customization/#custom-release-notes).

## Bump the version

In the future we will probably want to replace the hardcoded `Version` constant with [main.version](https://goreleaser.com/environment/#using-the-main-version). Here is a [stack overflow post on injecting to cmd/version.go](https://stackoverflow.com/a/47510909).

Commit this change on a branch along with the CHANGELOG updates in a single commit, and create a PR for merge to master.

## Cut a release

```bash
# Test run
goreleaser --skip-publish --snapshot --rm-dist

# Create a new tag on the master branch and push it
git tag -a v4.0.0 -m "Trying out GoReleaser"
git push origin v4.0.0

# Build and release
goreleaser --rm-dist
```

## Cut Release on GitHub

The generated archive files should be uploaded to the [draft release page created by GoReleaser](https://github.com/exercism/cli/releases). Describe the release (generally you can copy from the CHANGELOG.md file), select a specific commit to target, then test and publish the draft.
