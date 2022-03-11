# Changelog

Exercism configlet follows [semantic versioning](http://semver.org/).

----------------

## Next Release

* [**Your contribution here**](/CONTRIBUTING.md)

## v4.0.0 (30th October 2020)

* [#182](https://github.com/exercism/configlet/pull/182)
    Switch to GoReleaser [@ekingery]
* [#171](https://github.com/exercism/configlet/pull/171)
    Trim UUID in duplicate check [@ekingery]
* [#177](https://github.com/exercism/configlet/pull/177)
* [#167](https://github.com/exercism/configlet/pull/167)
* [#166](https://github.com/exercism/configlet/pull/166)
* [#164](https://github.com/exercism/configlet/pull/164)
* [#156](https://github.com/exercism/configlet/pull/156)
* [#151](https://github.com/exercism/configlet/pull/151)
    Various fetch-configlet bash script updates / fixes [@SaschaMann], [@SleeplessByte], [@ErikSchierboom], [@guygastineau], [@NobbZ]
* [#150](https://github.com/exercism/configlet/pull/150)
    Added --dry-run flag to the fmt command [@danielj-jordan]  

## v3.9.2 (11th June 2018)
* [#139](https://github.com/exercism/configlet/pull/139)
    Add missing boolean field ExerciseMetadata.AutoSupport [@kytrinyx]
* [#138](https://github.com/exercism/configlet/pull/138)
    Update type of ExerciseMetadata.ChecklistIssue to `int` [@kytrinyx]

## v3.9.1 (3rd June 2018)

* [#134](https://github.com/exercism/configlet/pull/134)
    Refactor the fmt command to rely on default marshaling of structs [@kytrinyx]
* [#133](https://github.com/exercism/configlet/pull/133)
    Add ordering to maintainers config [@kytrinyx]
* [#132](https://github.com/exercism/configlet/pull/132)
    Rework the fmt command and tests [@kytrinyx]
* [#131](https://github.com/exercism/configlet/pull/131)
    Distinguish between undefined and empty for topics [@kytrinyx]

## v3.9.0 (8th May 2018)

* [#127](https://github.com/exercism/configlet/pull/127)  Update configlet fmt to apply a definite ordering [@coriolinus]
* [#126](https://github.com/exercism/configlet/pull/126)  all: linting [@ferhatelmas]
* [#125](https://github.com/exercism/configlet/pull/125)  Gopkg: correct uuid dependency [@ferhatelmas]

## v3.8.0 (20th March 2018)

* [#122](https://github.com/exercism/configlet/pull/122) readme: Add documentation for the upgrade command - [@nywilken]
* [#119](https://github.com/exercism/configlet/pull/119) lint: Add checks for Invalid Unlocked By & Locked Core Exercises - [@tuxagon]
* [#118](https://github.com/exercism/configlet/pull/118) Add upgrade command - [@nywilken]
* [#115](https://github.com/exercism/configlet/pull/115) Add filter for unsupported exercise directories - [@nywilken]

## v3.7.0 (4th February 2018)

* [#98](https://github.com/exercism/configlet/pull/98) lint: Check for README presence - [@petertseng]
* [#111](https://github.com/exercism/configlet/pull/111) readme: Update UUID command example to new format - [@N-Parsons]
* [#110](https://github.com/exercism/configlet/pull/110) track: Check error on README generation - [@ferhatelmas]

## v3.6.1 (11th November 2017)

* [#106](https://github.com/exercism/configlet/pull/106) Replace broken UUID pkg - [@nywilken]
* [#105](https://github.com/exercism/configlet/pull/105) Git ignore .editorconfig - [@tleen]
* [#103](https://github.com/exercism/configlet/pull/103) tree: Fix invalid unlocked_by reference panic - [@tleen]
* [#104](https://github.com/exercism/configlet/pull/104) tree: fix bad indenting in fixture - [@tleen]

## v3.6.0 (20th October 2017)

* [#101](https://github.com/exercism/configlet/pull/101) Add UUID section to README - [@tleen]
* [#84](https://github.com/exercism/configlet/pull/84) Remove coupling between description and metadata - [@nywilken]
* [#97](https://github.com/exercism/configlet/pull/97) changelog: small typo - [@ferhatelmas]
* [#95](https://github.com/exercism/configlet/pull/95) Add generate command info to README - [@tleen]
* [#93](https://github.com/exercism/configlet/pull/93) Rename fixtures to match testing context - [@nywilken]
* [#90](https://github.com/exercism/configlet/pull/90) Add --track-id flag to lint command - [@robphoenix]
* [#83](https://github.com/exercism/configlet/pull/83) Add tree command - [@tleen]
* [#85](https://github.com/exercism/configlet/pull/85) Add support for custom readme titles - [@nywilken]
* [#88](https://github.com/exercism/configlet/pull/88) redirect output from ui package in lint test file - [@robphoenix]
* [#89](https://github.com/exercism/configlet/pull/89) Add a snake case version of the slug - [@kytrinyx]
* [#87](https://github.com/exercism/configlet/pull/87) Update topic format description - [@stkent]
* [#81](https://github.com/exercism/configlet/pull/81) Add links to authors in changelog - [@kytrinyx]

## v3.5.1 (23rd September 2017)

* [#78](https://github.com/exercism/configlet/pull/78) Fix error handling for malformed config files - [@nywilken]
* [#72](https://github.com/exercism/configlet/pull/72) Implement UI printer - [@robphoenix]
* [#74](https://github.com/exercism/configlet/pull/74) Add CHANGELOG.md - [@robphoenix]
* [#73](https://github.com/exercism/configlet/pull/73) Update CONTRIBUTING.md - [@rpottsoh]
* [#71](https://github.com/exercism/configlet/pull/71) Update the list of supported Go versions - [@nywilken]
* [#68](https://github.com/exercism/configlet/pull/68) Document the basic release process - [@kytrinyx]
* [#64](https://github.com/exercism/configlet/pull/64) Amendments to commands - [@robphoenix]

## v3.5.0 (5th September 2017)

* [#67](https://github.com/exercism/configlet/pull/67) Update track id detection to use absolute path - [@nywilken]
* [#63](https://github.com/exercism/configlet/pull/63) fmt: fix topic normalisation - [@robphoenix]
* [#57](https://github.com/exercism/configlet/pull/57) Add UUID validation across exercises and tracks to lint command - [@nywilken]

## v3.4.0 (22nd August 2017)

* [#59](https://github.com/exercism/configlet/pull/59) Replace calls to Fprintf with Fprintln - [@nywilken]
* [#55](https://github.com/exercism/configlet/pull/55) Refactor usage output for all commands - [@robphoenix]
* [#49](https://github.com/exercism/configlet/pull/49) Add check for missing problem-specifications directory - [@nywilken]
* [#52](https://github.com/exercism/configlet/pull/52) Add verbose flag to fmt - [@robphoenix]
* [#54](https://github.com/exercism/configlet/pull/54) Fix generate usage output - [@robphoenix]
* [#34](https://github.com/exercism/configlet/pull/34) Add fmt command - [@robphoenix]

## v3.3.0 (14th August 2017)

* [#48](https://github.com/exercism/configlet/pull/48) Provide a MixedCase version of the problem spec name - [@kytrinyx]
* [#47](https://github.com/exercism/configlet/pull/47) lint.go: Improve output of lint help command - [@robphoenix]
* [#45](https://github.com/exercism/configlet/pull/45) Minor changes to lint & generate commands - [@robphoenix]
* [#46](https://github.com/exercism/configlet/pull/46) Document how to run the tests on Windows - [@robphoenix]
* [#44](https://github.com/exercism/configlet/pull/44) Add CONTRIBUTING.md - [@robphoenix]
* [#43](https://github.com/exercism/configlet/pull/43) Vendor dependencies with dep - [@robphoenix]

## v3.2.0 (4th August 2017)

* [#37](https://github.com/exercism/configlet/pull/37) Lint maintainers - [@kytrinyx]
* [#38](https://github.com/exercism/configlet/pull/38) Fix golint complaints - [@kytrinyx]
* [#32](https://github.com/exercism/configlet/pull/32) configlet: Verify presence of test file - [@nywilken]
* [#36](https://github.com/exercism/configlet/pull/36) Delete unused variable - [@kytrinyx]
* [#31](https://github.com/exercism/configlet/pull/31) configlet: Update test lint command in doc.go - [@nywilken]

## v3.1.0 (23rd July 2017)

* [#26](https://github.com/exercism/configlet/pull/26) Add a command to generate UUIDs - [@kytrinyx]
* [#24](https://github.com/exercism/configlet/pull/24) Add generate command for exercise readmes - [@kytrinyx]
* [#23](https://github.com/exercism/configlet/pull/23) Rewrite with cobra - [@kytrinyx]
* [#25](https://github.com/exercism/configlet/pull/25) minor update readme.md - [@rpottsoh]

## v2.0.2 (9th July 2017)

## v2.0.1 (11th June 2017)

* [#19](https://github.com/exercism/configlet/pull/19) Make the error messages more expressive - [@kytrinyx]
* [#21](https://github.com/exercism/configlet/pull/21) Update license file - [@stkent]
* [#22](https://github.com/exercism/configlet/pull/22) Remove license info from README - [@stkent]
* [#20](https://github.com/exercism/configlet/pull/20) Update configlet check descriptions - [@stkent]

## v2.0.0 (10th June 2017)

* [#17](https://github.com/exercism/configlet/pull/17) Remove support for exercise implementations in the root of the track - [@kytrinyx]
* [#16](https://github.com/exercism/configlet/pull/16) Remove support for problems key in config.json - [@kytrinyx]
* [#14](https://github.com/exercism/configlet/pull/14) Readme.md: Link to config.json documentation. - [@Insti]

## v1.1.0 (8th May 2017)

* [#13](https://github.com/exercism/configlet/pull/13) Configurable example solution pattern - [@Insti]
* [#9](https://github.com/exercism/configlet/pull/9) Simplify gofmt style - [@ferhatelmas]
* [#6](https://github.com/exercism/configlet/pull/6) Run tests on travis - [@kytrinyx]
* [#5](https://github.com/exercism/configlet/pull/5) Ignore img directory - [@ryanplusplus]

[@coriolinus]: https://github.com/coriolinus
[@Insti]: https://github.com/Insti
[@ferhatelmas]: https://github.com/ferhatelmas
[@kytrinyx]: https://github.com/kytrinyx
[@N-Parsons]: https://github.com/N-Parsons
[@nywilken]: https://github.com/nywilken
[@robphoenix]: https://github.com/robphoenix
[@rpottsoh]: https://github.com/rpottsoh
[@ryanplusplus]: https://github.com/ryanplusplus
[@petertseng]: https://github.com/petertseng
[@stkent]: https://github.com/stkent
[@tleen]: https://github.com/tleen
[@tuxagon]: https://github.com/tuxagon
[@danielj-jordan]: https://github.com/danielj-jordan
[@ekingery]: https://github.com/ekingery
[@ErikSchierboom]: https://github.com/ErikSchierboom
[@SleeplessByte]: https://github.com/SleeplessByte
[@SaschaMann]: https://github.com/SaschaMann
[@guygastineau]: https://github.com/guygastineau
[@NobbZ]: https://github.com/NobbZ
