# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## [0.0.11] - 2021-03-19
### Added
- flag `--titlePrettify` to apply strings.Title in message title and remove these characters -, /, \

### Changed
- update golang to 1.15

## [0.0.10] - 2021-12-21
### Added
- flag to parse annotations using strings.HasSuffix to create button: -S or --annotations-suffix-link
- flag to exclude annotations when creating buttons: -E or --annotations-suffix-exclude

## [0.0.9] - 2020-11-30
### Added
- Annotations and labels parsing checking

## [0.0.8] - 2020-11-30
### Added
- Add `--annotations-as-link` flag to parse annotation as link

## [0.0.7] - 2020-11-20
### Changed
- Fix bonsai asset list


## [0.0.6] - 2020-11-20
### Changed
- Changed go to 1.14

### Removed
- Removed travis-ci integration

### Added
- Added tests
- Added golanglint


## [0.0.5] - 2020-11-20
### Changed
- Changed webhook URL environment variable from `WEBHOOK_HANGOUTSCHAT` to `HANGOUTSCHAT_WEBHOOK`.
- Changed from spf13/cobra to sensu-community/sensu-plugin-sdk. 
- Changed `--withAnnotations` to parse all annotations, and exclude if it contains `sensu.io/plugins/sensu-hangouts-chat-handler/config`, and send as text to Hangouts Chat. 

### Added
- Added `--withLabels` to parse all labels, and exclude if it contains `sensu.io/plugins/sensu-hangouts-chat-handler/config`, and send as text to Hangouts Chat.

### Removed 
- Removed annotations as link in message.

## [0.0.4] - 2020-02-25
### Changed
- Change from text message to card message in Hangouts Chat

## [0.0.3] - 2020-01-16
### Changed
- Change from dep to go mod.
- gometalinter to golangci-lint
- correct goreleaser

## [0.0.2] - 2019-11-24
### Added
- Add `HANGOUTS_ANNOTATIONS` to parse annotations to include that information inside the alert.
### Changed
- Update README.


## [0.0.1] - 2019-10-30

### Added
- Initial release