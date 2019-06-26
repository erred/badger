# badger

badges for GCP cloud build

[![License](https://img.shields.io/github/license/seankhliao/badger.svg?style=for-the-badge&maxAge=31536000)](LICENSE)
[![Build](https://badger.seankhliao.com/i/github_seankhliao_badger)](https://badger.seankhliao.com/l/github_seankhliao_badger)

## About

Cloud build doesn't natively support badges yet :(

#### Run

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
