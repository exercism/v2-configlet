# Cutting a Configlet Release

`configlet` uses [GoReleaser](https://goreleaser.com) to automate the
release process. 

## Requirements

1. [Install GoReleaser](https://goreleaser.com/install/)
1. [Setup GitHub token](https://goreleaser.com/environment/#github-token)
1. Have a gpg key installed on your machine - it is [used for signing the artifacts](https://goreleaser.com/customization/sign/)

## Confirm / Update the Changelog

Make sure all the recent changes are reflected in the "next release" section of the CHANGELOG.md file.  All the changes in the "next release" section should be moved to a new section that describes the version number, and gives it a date.

You can view changes using the /compare/ view:
https://github.com/exercism/configlet/compare/$PREVIOUS_RELEASE...master

Some features to consider - GoReleaser supports the [auto generation of a changelog](https://goreleaser.com/customization/#customize-the-changelog) and has a [release notes feature](https://goreleaser.com/customization/#custom-release-notes).


## Cut a release

Create a PR for the updates to the CHANGELOG. Once that is merged to master, you should test the builds by running: `goreleaser --skip-publish --snapshot --rm-dist`

Once you have verified that the binaries are built properly:

```bash
# Create a new tag on the master branch and push it
git tag -a v3.9.10 -m "Trying out GoReleaser"
git push origin v3.9.10

# Build and release
goreleaser --rm-dist
```

### GitHub Follow-Up / Backwards Compatibility

The `fetch-configlet` file is hardcoded into most track repos, so we need to rename the `*.tar.gz` files to `*.tgz` and upload these manually to the release for backwards compatibility. This command will do the appropriate renaming: `rename 's/\.tar.gz\z/.tgz/' dist/*`

Upload the renamed files to the [draft release page created by GoReleaser](https://github.com/exercism/configlet/releases). Describe the release (generally you can copy from the CHANGELOG.md file), select a specific commit to target, then test and publish the draft. Make sure to run the [fetch_configlet](scripts/fetch-configlet) script to confirm functionality once the draft is published.
