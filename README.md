# badger

[![Build](https://badger.seankhliao.com/i/github_seankhliao_badger)](https://badger.seankhliao.com/l/github_seankhliao_badger)
[![License](https://img.shields.io/github/license/seankhliao/badger.svg?style=for-the-badge)](LICENSE)

badges for GCP Cloud Build

## Usage

```
badger [-p 8080] [-pr com-seankhliao]
  -p    port to listen on
  -pr   GCP project to query
```

Accepts the following urls:

> github repos are the form `github_$USER_$REPO`

- `/i/$REPO`: redirect to the appropriate shields.io img
- `/l/$REPO`: redirect to the build history in GCP console

ex:

```
[![Build][build-img]][build-link]

[build-img]: https://badger.seankhliao.com/i/github_seankhliao_badger
[build-link]: https://badger.seankhliao.com/l/github_seankhliao_badger
```

## Design

steps:

1. Query GCP for builds list fitting repo filter
2. Filter skipping working / queued / cancelled, use first result
3. generate JSON

## TODO

- manage caching
- doc url
- markdown generator
