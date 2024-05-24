<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry is required to include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The tag should consist of where the change is being made ex. (x/staking), (store)
The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
"CLI Breaking" for breaking CLI commands.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

Here is the corrected version of the changelog with proper grammar and typo corrections:

# Changelog

## [Unreleased]

### Features

- (x/token-convert) [#18](https://github.com/tabilabs/tabi/pull/18) Token-convert module implements the token conversion between Tabi and VeTabi.
- (x/captains) [#20](https://github.com/tabilabs/tabi/pull/20) Captains module manages captain nodes and calculates rewards for captain nodes in epochs.
- (x/claims) [#23](https://github.com/tabilabs/tabi/pull/23) Claims module allows users to claim rewards accrued from the Captains module.
- (x/mint) [#31](https://github.com/tabilabs/tabi/pull/31) Mint module mints inflation rewards for network participants.

### Improvements

- (ante) [#22](https://github.com/tabilabs/tabi/pull/28) Add an allowlist ante handler restricting EVM transactions.

